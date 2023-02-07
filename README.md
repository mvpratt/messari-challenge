# messari-challenge

## Requirements:
1) Read a series of JSON objects from stdin

Example input format:

```
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

```
{
    "market":5775,
	  "total_volume":1234567.89,
	  "mean_price":23.33,
	  "mean_volume":6144.299,
	  "volume_weighted_average_price":5234.2,
	  "percentage_buy":0.50
}
```

## Performance:
Typical: ~16 seconds on M1 Macbook Pro