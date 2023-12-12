# Livepeer_value_redeemed publisher

A simple Prometheus exporter that can be used to temporarily increment the `livepeer_value_redeemed` vector counter metric with a given value and publish it on http://localhost:7935/metrics. This can be used to correct the `livepeer_value_redeemed` metric due to a bug that was present in the `go-livepeer` binary before [go-livepeer/2916](https://github.com/livepeer/go-livepeer/pull/2916) was merged.

## Usage

### Build

```bash
go build
```

### Run

1. Run the [go-livepeer](https://github.com/livepeer/go-livepeer) binary without the `-monitor` flag.
2. To increment the metric, Run the `livepeer_value_redeemed` binary with the `-value` and `-node_id` flags. The value should be in wei and represent the error in the `livepeer_value_redeemed` metric.

   ```bash
   ./livepeer_value_redeemed -value number -node_id string
   ```

3. Wait for the value to be published.
4. Stop the `livepeer_value_redeemed` binary.
5. Start the [go-livepeer](https://github.com/livepeer/go-livepeer) binary with the `-monitor` flag.

> [!NOTE]\
> You can also change the `prometheus.yml` to listen for Livepeer metrics on a different port than the default `7935` port and then run the `livepeer_value_redeemed` binary with the `-port` flag.
