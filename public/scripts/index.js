const sseEndpoints = {
    humidity: '/api/v1/humidity/event',
    mass: '/api/v1/mass/event',
    ammonia: '/api/v1/ammonia/event',
    temperature: '/api/v1/temperature/event'
};

// Define maximum values for each event type
const maxValues = {
    humidity: 100,
    mass: 5,
    ammonia: 4095,
    temperature: 80
};

const units = {
    humidity: '%',
    mass: 'kg',
    ammonia: 'NH3',
    temperature: '°C'
};

const gaugeContainer = d3.select("#gauge-container");

// Create gauge plots
const width = 200;
const height = 300; // Adjusted height to fit the gauge and labels
const minValue = 0;
const maxAngle = Math.PI;

for (const eventType in sseEndpoints) {
    const svg = gaugeContainer.append("svg")
        .attr("width", width)
        .attr("height", height); // Adjusted height

    const gauge = svg.append("g")
        .attr("transform", `translate(${width / 2},${height / 1.5})`); // Adjusted position

    const arc = d3.arc()
        .innerRadius(80)
        .outerRadius(100)
        .startAngle(-Math.PI / 2);

    const background = gauge.append("path")
        .datum({ endAngle: Math.PI / 2 })
        .attr("class", "background")
        .attr("d", arc)
        .attr("fill", "#e0e0e0");

    const foreground = gauge.append("path")
        .datum({ endAngle: -Math.PI / 2 })
        .attr("class", "foreground")
        .attr("d", arc)
        .attr("fill", "#009688");

    const valueText = gauge.append("text")
        .attr("class", "value")
        .attr("dy", "0.35em")
        .attr("text-anchor", "middle")
        .attr("font-size", "18px");

    // Add label for the event type
    svg.append("text")
        .attr("class", "event-label")
        .attr("x", width / 2)
        .attr("y", height - 10)
        .attr("text-anchor", "middle")
        .text(eventType.charAt(0).toUpperCase() + eventType.slice(1)); // Capitalize first letter

    // Add unit label
    svg.append("text")
        .attr("class", "unit-label")
        .attr("x", width / 2)
        .attr("y", height - 30)
        .attr("text-anchor", "middle")
        .text(units[eventType]);

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

            valueText.text(value.toFixed(2));
        }
    };

    sse.onerror = function(event) {
        console.error(`SSE Error (${eventType}):`, event);
    };
}