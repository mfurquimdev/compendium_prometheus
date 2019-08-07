package br.com.furquim.servlets;

import io.prometheus.client.Counter;
import io.prometheus.client.Histogram;
import io.prometheus.client.exporter.MetricsServlet;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.concurrent.TimeUnit;

import static java.lang.Math.random;


@WebServlet(name = "TestServlet", urlPatterns = "/test")
public class TestServlet extends HttpServlet {
    static final Counter requests = Counter.build()
            .name("requests_total").help("Total número de requisições.").register();

    static final Histogram histRand = Histogram.build()
            .buckets(0.1, 0.25, 0.5, 0.75, 0.9, 1)
            .name("requests_random_numbers").help("Random number generated").register();

    static final Histogram latency = Histogram.build()
            .name("requests_seconds").help("Latência das requisições").register();

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
        requests.inc();
        histRand.observe(Math.random());
        System.out.println("TestServlet.doGet");
        resp.getWriter().println("Número de requisições: " + requests.get());

        resp.getWriter().println("Histogram Random:");
        for (int i = 0; i < histRand.collect().size(); i++) {
            for (int j = 0; j < histRand.collect().get(i).samples.size(); j++) {
                resp.getWriter().println("\t" + histRand.collect().get(i).samples.get(j));
            }
        }

        resp.getWriter().println("Histogram Latency:");
        for (int i = 0; i < latency.collect().size(); i++) {
            for (int j = 0; j < latency.collect().get(i).samples.size(); j++) {
                resp.getWriter().println("\t" + latency.collect().get(i).samples.get(j));
            }
        }
    }
}
