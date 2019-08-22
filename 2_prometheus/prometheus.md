Prometheus has a ticker which triggers the scrape exaclty on \<scrape\_interval\> time, regardless of the time it took to scrape the target. [1](https://utcc.utoronto.ca/~cks/space/blog/sysadmin/PrometheusScrapeIntervalBit)


it is recommended that \<scrape\_interval\> and \<evaluation\_interval\> are equal. [2](https://stackoverflow.com/questions/52167869/scrape-interval-and-evaluation-interval-in-prometheus)

\<scrapetimeout\> is set for each job


