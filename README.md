# messari-challenge

## Reference
https://engineering.messari.io/blog/messari-market-data-coding-challenge

## High Level Requirements
1) Read a series of JSON objects from stdin

Example input format:

```json
BEGIN
{"id":1,"market":1,"price":1.8126489025824468,"volume":2529.667281360351,"is_buy":true}
{"id":2,"market":2,"price":2.1663930707558356,"volume":3370.4751940246724,"is_buy":false}
{"id":3,"market":11,"price":11.812182638400632,"volume":1644.641438186002,"is_buy":true}
END
Trade Count:  10
```

2) Send to stdout the following metrics for each market, also in JSON.
One resulting object per market.

Example output:

```json
{
    "market":5775,
	  "total_volume":1234567.89,
	  "mean_price":23.33,
	  "mean_volume":6144.299,
	  "volume_weighted_average_price":5234.2,
	  "percentage_buy":0.50
}
```

3. Complete this processing as fast as possible

## Performance
Typical: ~16 seconds on M1 Macbook Pro


## Detailed Prompt

### Objective
Build an efficient tool to compute aggregate market data from raw trades.

### Background
Among many other things, Messari is tasked with keeping track of asset prices, volume figures, etc. for as many crypto assets and crypto markets as possible. Fundamentally, the tasks at hand require a high-performance codebase. This exercise is designed to assess how you handle a problem that makes you think about low-level performance optimizations in the context of crypto trade data.

### Details
As part of this exercise prompt, you will be provided with a binary that writes ten million trade objects as JSON to stdout, and then exits. It writes one object per line, as quickly as possible. Your objective is to write a tool that parses each trade as it comes in and computes various aggregate metrics from the provided data, completing the set of ten million trades in as little time as possible.

Your tool should accept input from stdin so output from the provided binary can be piped into your tool using a terminal window.

As an example, three consecutive lines of output would be of the form below:

```json
{"id":121509,"market":5773,"price":1.234,"volume":1234.56,"is_buy":true} {"id":121510,"market":5774,"price":2.345,"volume":2345.67,"is_buy":false} {"id":121511,"market":5775,"price":3.456,"volume":3456.78,"is_buy":true}
```

Where:

ID is a unique auto-incrementing int starting at 1. You can use it to keep track of how many trades you have processed so far.
Market is a random int between 1 and about 12,000. Each trade belongs to a market.
Price and volume are floats to varying degrees of precision
IsBuy is a bool denoting whether the trade is a buy or a sell
Your app should continuously compute and keep track of the following as new trades come in:

Total volume per market
Mean price per market
Mean volume per market
Volume-weighted average price per market
Percentage buy orders per market
Note: The trade objects provided in this project do not have timestamps, so all the above should be calculated across all the trades observed for that market across the given set of ten million.
Once all trades have been parsed, your app should send to stdout the following metrics for each market, also in JSON. One resulting object per market. For instance, you might output:

```json
{"market":5775,"total_volume":1234567.89,"mean_price":23.33,"mean_volume":6144.299,"volume_weighted_average_price":5234.2,"percentage_buy":0.50}
```

(But don’t worry about copying these exact keys or anything, feel free to represent the percent as out of 1 or out of 100, etc.)

### Details on the Provided Binary Behavior
When started, the provided binary will write a line containing “BEGIN”, followed by about ten million lines each containing a JSON object as indicated in the previous example. Then it will write a line containing “END” and provide simple stats on what has just occurred. These stats are the number of trades written, the total number of markets included in the run (which can vary, but be up to about 12,000), and the amount of time the binary was running.

You only need to pay attention to JSON after “BEGIN” and before “END”. Feel free to only pay attention to valid JSON and ignore the other lines—they are included just in case you want to pass them through and inspect the market count and runtime after each run.
To observe the binary behavior for yourself, you could execute the following, for instance, on an arm64 MacBook:

```
> chmod +x ./stdoutinator_arm64_darwin.bin > ./stdoutinator_arm64_darwin.bin > ./output.txt
```

and then look at output.txt in a text editor.

We have also included the source, so you can look over it, and so you can compile it yourself if you prefer.

### Details on How Your Work Will Be Assessed
Your work will be assessed primarily on how performant your resulting project is. There is no specific execution time requirement, but it should be sufficiently performant according to your own judgment.