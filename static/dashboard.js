// dashboard.js: Periodically fetch metrics and update the table

function formatDuration(s) {
    if (s < 1) return (s * 1000).toFixed(0) + ' ms';
    if (s < 60) return s.toFixed(2) + ' s';
    return (s / 60).toFixed(2) + ' s';
}
function formatMB(b) {
    return (b/1024/1024).toFixed(2) + ' MB';
}

function updateMetricsTable() {
    fetch('/api/metrics')
        .then(resp => resp.json())
        .then(functions => {
            const tbody = document.querySelector('.metrics-table tbody');
            tbody.innerHTML = '';
            functions.forEach(fn => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td><a href="/function?name=${fn.name}">${fn.name}</a></td>
                    <td>${fn.duration ? formatDuration(fn.duration) : ''}</td>
                    <td>${fn.avg_duration ? formatDuration(fn.avg_duration) : ''}</td>
                    <td>${fn.mem_alloc ? formatMB(fn.mem_alloc) : ''}</td>
                    <td>${fn.invocations}</td>
                `;
                tbody.appendChild(tr);
            });
        });
}

setInterval(updateMetricsTable, 1000);
document.addEventListener('DOMContentLoaded', updateMetricsTable);
