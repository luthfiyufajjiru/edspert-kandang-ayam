const sseEndpoints = {
    humidity: '/api/v1/humidity/event',
    mass: '/api/v1/mass/event',
    ammonia: '/api/v1/ammonia/event',
    temperature: '/api/v1/temperature/event'
};

// Define maximum values for each event type
const maxValues = {
    humidity: 100,
    mass: 10,
    ammonia: 4500,
    temperature: 120
};

const gaugeContainer = d3.select("#gauge-container");

// Create gauge plots
const width = 200;
const height = 200;
const minValue = 0;

for (const eventType in sseEndpoints) {
    const svg = gaugeContainer.append("svg")
        .attr("width", width)
        .attr("height", height);

    // Add label for the gauge
    const labelMap = {
        mass: "Mass (kg)",
        ammonia: "Ammonia (NH3)",
        temperature: "Temperature (°C)",
        humidity: "Humidity (%)"
    };

    svg.append("text")
        .attr("class", "label")
        .attr("x", width / 2)
        .attr("y", height - 10)
        .attr("text-anchor", "middle")
        .text(labelMap[eventType]);

    const gauge = svg.append("g")
        .attr("transform", `translate(${width / 2},${height / 2})`);

    const arc = d3.arc()
        .innerRadius(40)
        .outerRadius(60)
        .startAngle(-Math.PI / 2);

    const background = gauge.append("path")
        .datum({ endAngle: Math.PI / 2 })
        .attr("class", "background")
        .attr("d", arc)
        .attr("stroke", "gray")
        .attr("stroke-width", 10);

    const foreground = gauge.append("path")
        .datum({ endAngle: -Math.PI / 2 })
        .attr("class", "foreground")
        .attr("d", arc)
        .attr("stroke", "blue")
        .attr("stroke-width", 10);

    const valueText = gauge.append("text")
        .attr("class", "value")
        .attr("dy", "-0.5em")
        .attr("text-anchor", "middle");

    // SSE Connection
    const sse = new EventSource(sseEndpoints[eventType]);

    sse.onmessage = function(event) {
        const eventData = JSON.parse(event.data);
        const data = eventData.data;

        if (Array.isArray(data) && data.length > 0) {
            const value = parseFloat(data[data.length - 1]);

            // Update gauge with new value
            const angle = d3.scaleLinear()
                .domain([minValue, maxValues[eventType]])
                .range([-Math.PI / 2, Math.PI / 2]);

            foreground.transition()
                .duration(1000)
                .attrTween("d", function(d) {
                    const interpolate = d3.interpolate(d.endAngle, angle(value));
                    return function(t) {
                        d.endAngle = interpolate(t);
                        return arc(d);
                    };
                });

            valueText.text(value);
        }
    };

    sse.onerror = function(event) {
        console.error(`SSE Error (${eventType}):`, event);
    };
}