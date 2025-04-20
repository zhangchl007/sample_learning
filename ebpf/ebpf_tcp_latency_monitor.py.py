#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Ensure bcc is installed: pip install bcc
from bcc import BPF
import time
import socket
import struct
import math  # Use math for sqrt
import datetime
import argparse

# BPF 程序(C 语言部分）
bpf_text = """
#include <uapi/linux/ptrace.h>
#include <net/sock.h>
#include <net/tcp_states.h>
#include <bcc/proto.h>
#include <linux/tcp.h>
#include <linux/ip.h>

// Key for tracking in-flight packets
struct inflight_key {
    u32 saddr;
    u32 daddr;
    u16 sport;
    u16 dport;
    u32 seq;
};

// Value: Timestamp when the packet was sent
struct inflight_val {
    u64 send_ts_ns;
};

// Key for aggregated statistics per connection
struct stats_key {
    u32 saddr;
    u32 daddr;
    u16 sport;
    u16 dport;
};

// Value: Aggregated RTT statistics
struct stats_val {
    u64 rtt_ns_sum;
    u64 rtt_ns_sum_sq;
    u64 rtt_ns_min;
    u64 rtt_ns_max;
    u64 count;
};

// Map to store timestamps of in-flight data segments
BPF_HASH(inflight_map, struct inflight_key, struct inflight_val, 10240);

// Map to store aggregated statistics
BPF_HASH(stats_map, struct stats_key, struct stats_val, 10240);

// Map to store retransmission counts
BPF_HASH(retrans_map, struct stats_key, u64, 10240);

// Kprobe on tcp_transmit_skb: Called when TCP prepares to send a segment
int trace_transmit(struct pt_regs *ctx) {
    struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
    u16 family = sk->__sk_common.skc_family;
    if (family != AF_INET) {
        return 0; // Only IPv4 for simplicity
    }

    u32 saddr = sk->__sk_common.skc_rcv_saddr;
    u32 daddr = sk->__sk_common.skc_daddr;
    u16 sport = sk->__sk_common.skc_num;
    u16 dport = sk->__sk_common.skc_dport;
    dport = ntohs(dport);

    struct tcp_sock *tp = (struct tcp_sock *)sk;
    u32 ack_seq = tp->rcv_nxt; // This acknowledges data *up to* ack_seq - 1

    // Create lookup key
    struct inflight_key key = {};
    key.saddr = saddr;
    key.daddr = daddr;
    key.sport = sport;
    key.dport = dport;
    key.seq = ack_seq - 1;

    struct inflight_val val = {};
    val.send_ts_ns = bpf_ktime_get_ns();

    inflight_map.update(&key, &val);

    return 0;
}

// Kprobe on tcp_ack: Called when processing a received ACK
int trace_ack(struct pt_regs *ctx, struct sock *sk) {
    u16 family = sk->__sk_common.skc_family;
    if (family != AF_INET) {
        return 0;
    }

    u32 saddr = sk->__sk_common.skc_rcv_saddr;
    u32 daddr = sk->__sk_common.skc_daddr;
    u16 sport = sk->__sk_common.skc_num;
    u16 dport = sk->__sk_common.skc_dport;
    dport = ntohs(dport);

    struct tcp_sock *tp = (struct tcp_sock *)sk;
    u32 ack_seq = tp->rcv_nxt;

    // Create lookup key
    struct inflight_key key = {};
    key.saddr = saddr;
    key.daddr = daddr;
    key.sport = sport;
    key.dport = dport;
    key.seq = ack_seq-1;

    struct inflight_val *valp = inflight_map.lookup(&key);

    if (valp) {
        u64 send_ts = valp->send_ts_ns;
        u64 ack_ts = bpf_ktime_get_ns();
        inflight_map.delete(&key); // Remove entry once ACKed

        if (ack_ts <= send_ts) {
            return 0; // Invalid timestamp
        }

        u64 rtt_ns = ack_ts - send_ts;

        // Update aggregate statistics
        struct stats_key stats_key = {};
        stats_key.saddr = saddr;
        stats_key.daddr = daddr;
        stats_key.sport = sport;
        stats_key.dport = dport;

        struct stats_val zero = {};
        zero.rtt_ns_min = (u64)-1;

        struct stats_val *stats = stats_map.lookup_or_init(&stats_key, &zero);

        if (stats) {
            stats->rtt_ns_sum += rtt_ns;
            stats->rtt_ns_sum_sq += rtt_ns * rtt_ns;
            if (rtt_ns < stats->rtt_ns_min) stats->rtt_ns_min = rtt_ns;
            if (rtt_ns > stats->rtt_ns_max) stats->rtt_ns_max = rtt_ns;
            stats->count += 1;
        }
    }
    return 0;
}

// Kprobe on tcp_retransmit_skb: Count retransmissions
int trace_retransmit(struct pt_regs *ctx, struct sock *sk) {
    u16 family = sk->__sk_common.skc_family;
    if (family != AF_INET) return 0;
    struct stats_key key = {};
    key.saddr = sk->__sk_common.skc_rcv_saddr;
    key.daddr = sk->__sk_common.skc_daddr;
    key.sport = sk->__sk_common.skc_num;
    key.dport = ntohs(sk->__sk_common.skc_dport);
    u64 zero = 0, *cnt;
    cnt = retrans_map.lookup_or_init(&key, &zero);
    if (cnt) (*cnt)++;
    return 0;
}
"""

# Initialize BPF
b = BPF(text=bpf_text)
# Load BPF program
try:
    b.attach_kprobe(event="__tcp_transmit_skb", fn_name="trace_transmit")
    b.attach_kprobe(event="tcp_ack", fn_name="trace_ack")
    b.attach_kprobe(event="tcp_retransmit_skb", fn_name="trace_retransmit")
except Exception as e:
    print("Failed to attach kprobes:", e)
    exit(1)

print("Tracing TCP RTT (per-packet)... Ctrl-C to stop.")

def ip_to_str(ip):
    return socket.inet_ntop(socket.AF_INET, struct.pack("I", ip))

# Function to get port from network byte order u16
def port_to_str(port):
     # Assuming port is already in host byte order if read directly from key
     # If port was network byte order in key, use socket.ntohs(port)
     return str(port)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="eBPF TCP RTT Monitor")
    parser.add_argument("--dest-ip", type=str, default="142.250.189.4", help="Destination IP to filter (default: 142.250.189.4)")
    args = parser.parse_args()
    dest_ip = args.dest_ip

    try:
        last_reset = time.time()
        window_sec = 600  # 10 minutes

        while True:
            time.sleep(10)
            now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            print("\n[%s] %-19s:%-5s -> %-15s:%-5s | avg min max stddev(ms) | count" %
                  (now, "SOURCE", "SPORT", "DEST", "DPORT"))
            # ...existing code inside while True: ...
            stats_map = b.get_table("stats_map")
            retrans_map = b.get_table("retrans_map")
            print("stats_map size:", len(stats_map))
            for k, v in stats_map.items():
                saddr = ip_to_str(k.saddr)
                daddr = ip_to_str(k.daddr)
                if daddr != dest_ip:
                    continue  # Only show stats for the desired destination

                sport = port_to_str(k.sport)
                dport = port_to_str(k.dport)

                avg_rtt_ms = 0
                min_rtt_ms = 0
                max_rtt_ms = 0
                stddev_rtt_ms = 0

                if v.count > 0:
                    avg_rtt_ns = v.rtt_ns_sum / v.count
                    avg_rtt_ms = avg_rtt_ns / 1e6
                    max_rtt_ms = v.rtt_ns_max / 1e6
                    min_rtt_ms = 0 if v.rtt_ns_min == (2**64 - 1) else v.rtt_ns_min / 1e6

                    if v.count > 1:
                        variance_ns = (v.rtt_ns_sum_sq / v.count) - (avg_rtt_ns * avg_rtt_ns)
                        variance_ns = max(0, variance_ns)
                        stddev_ns = math.sqrt(variance_ns)
                        stddev_rtt_ms = stddev_ns / 1e6

                # Get retransmission count for this connection
                retrans_key = k
                retrans_count = retrans_map.get(retrans_key)
                retrans_count_val = retrans_count.value if retrans_count else 0

                print("%-19s:%-5s -> %-15s:%-5s | %6.2f %6.2f %6.2f %6.2f | %6d | retrans: %d" %
                      (saddr, sport, daddr, dport,
                       avg_rtt_ms, min_rtt_ms, max_rtt_ms, stddev_rtt_ms, v.count, retrans_count_val))

    except KeyboardInterrupt:
        print("Exiting...")