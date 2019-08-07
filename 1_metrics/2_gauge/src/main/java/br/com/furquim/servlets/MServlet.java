package br.com.furquim.servlets;

import io.prometheus.client.exporter.MetricsServlet;

import javax.servlet.annotation.WebServlet;

@WebServlet(name = "MServlet", urlPatterns = "/metrics")
public class MServlet extends MetricsServlet {


}
