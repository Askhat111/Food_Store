async function loadReports() {
    loadLowStockReport();
    loadExpiringReport();
    loadDailySalesReport();
}

async function loadLowStockReport() {
    try {
        const response = await fetch(`${API_URL}/reports/low-stock`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const products = await response.json();
        
        const container = document.getElementById('lowStockReport');
        if (products && products.length > 0) {
            container.innerHTML = products.map(p => `
                <div class="report-item">
                    <div class="report-item-name">${p.name}</div>
                    <div class="report-item-value">
                        <span class="status-badge status-low">${p.quantity} / ${p.min_stock_threshold || 5}</span>
                    </div>
                </div>
            `).join('');
        } else {
            container.innerHTML = '<div class="report-empty">No low stock items</div>';
        }
    } catch (error) {
        console.error('Error loading low stock report:', error);
    }
}

async function loadExpiringReport() {
    try {
        const response = await fetch(`${API_URL}/reports/expiring?days=7`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const products = await response.json();
        
        const container = document.getElementById('expiringReport');
        if (products && products.length > 0) {
            container.innerHTML = products.map(p => {
                let daysText = '';
                if (p.days_to_expiry < 0) {
                    daysText = `Expired ${-p.days_to_expiry}d ago`;
                } else if (p.days_to_expiry === 0) {
                    daysText = 'Expires today';
                } else {
                    daysText = `${p.days_to_expiry} days left`;
                }
                
                return `
                    <div class="report-item">
                        <div class="report-item-name">${p.name}</div>
                        <div class="report-item-value">
                            <span class="status-badge status-expiring">${daysText}</span>
                        </div>
                    </div>
                `;
            }).join('');
        } else {
            container.innerHTML = '<div class="report-empty">No expiring products</div>';
        }
    } catch (error) {
        console.error('Error loading expiring report:', error);
    }
}

async function loadDailySalesReport() {
    try {
        const response = await fetch(`${API_URL}/reports/daily-sales`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const sales = await response.json();
        
        const container = document.getElementById('dailySalesReport');
        if (sales && sales.length > 0) {
            let totalSales = 0;
            let totalTransactions = 0;
            
            sales.forEach(s => {
                totalSales += s.total;
                totalTransactions += s.count;
            });
            
            container.innerHTML = `
                <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 15px; margin-bottom: 20px;">
                    <div class="report-item">
                        <div class="report-item-name">Total Sales</div>
                        <div class="report-item-value" style="font-size: 1.5em; color: #27ae60; font-weight: 600;">
                            ${totalSales.toFixed(2)} ₸
                        </div>
                    </div>
                    <div class="report-item">
                        <div class="report-item-name">Transactions</div>
                        <div class="report-item-value" style="font-size: 1.5em; color: #3498db; font-weight: 600;">
                            ${totalTransactions}
                        </div>
                    </div>
                    <div class="report-item">
                        <div class="report-item-name">Avg per Transaction</div>
                        <div class="report-item-value" style="font-size: 1.5em; color: #9b59b6; font-weight: 600;">
                            ${(totalSales / totalTransactions).toFixed(2)} ₸
                        </div>
                    </div>
                </div>
                <table>
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Sales (₸)</th>
                            <th>Transactions</th>
                            <th>Avg (₸)</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${sales.map(s => `
                            <tr>
                                <td>${s.date}</td>
                                <td>${s.total.toFixed(2)} ₸</td>
                                <td>${s.count}</td>
                                <td>${s.avg_sale.toFixed(2)} ₸</td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
        } else {
            container.innerHTML = '<div class="report-empty">No sales data available</div>';
        }
    } catch (error) {
        console.error('Error loading daily sales report:', error);
    }
}