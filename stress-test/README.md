# Stress Test Results

The images in each folder contain the gathered metrics when running Fn Server under stress by using gattling.
The conclusion we can extract is that the scaling system works quite well on Fn which enables it to handle quite
a lot of stress, as much as 500 requests per second when doing simple tasks.

When trying to stress test Fn Flow we encounter the problem that it does not scale at all, and it dies as soon as it
gets near a hundred requests per second, but having a really bad performance and latency even with relatively low requests.
