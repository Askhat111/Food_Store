async function loadSales() {
    try {
        const response = await fetch(`${API_URL}/sales`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const sales = await response.json();
        
        const tbody = document.getElementById('salesTableBody');
        tbody.innerHTML = '';
        
        sales.slice(0, 20).forEach(sale => {
            const tr = document.createElement('tr');
            const date = new Date(sale.sale_date).toLocaleString();
            tr.innerHTML = `
                <td>${date}</td>
                <td>${sale.product_name}</td>
                <td>${sale.quantity}</td>
                <td>${sale.total.toFixed(2)} ₸</td>
                <td>${sale.user_name}</td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error('Error loading sales:', error);
    }
}

async function makeSale(event) {
    event.preventDefault();
    
    const sale = {
        product_id: parseInt(document.getElementById('saleProductId').value),
        quantity: parseInt(document.getElementById('saleQuantity').value)
    };
    
    try {
        const response = await fetch(`${API_URL}/sales`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(sale)
        });
        
        if (response.ok) {
            const data = await response.json();
            alert(`Sale recorded successfully!\nTotal: ${data.total.toFixed(2)} ₸`);
            document.getElementById('saleForm').reset();
            loadSales();
            loadProducts();
            loadAlerts();
        } else {
            const error = await response.text();
            alert('Error: ' + error);
        }
    } catch (error) {
        alert('Error recording sale: ' + error.message);
    }
}

async function loadPurchases() {
    try {
        const response = await fetch(`${API_URL}/purchases`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const purchases = await response.json();
        
        const tbody = document.getElementById('purchasesTableBody');
        tbody.innerHTML = '';
        
        purchases.slice(0, 20).forEach(purchase => {
            const tr = document.createElement('tr');
            const date = new Date(purchase.purchase_date).toLocaleString();
            tr.innerHTML = `
                <td>${date}</td>
                <td>${purchase.product_name}</td>
                <td>${purchase.quantity}</td>
                <td>${purchase.price.toFixed(2)} ₸</td>
                <td>${purchase.total.toFixed(2)} ₸</td>
                <td>${purchase.supplier || '-'}</td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error('Error loading purchases:', error);
    }
}

async function makePurchase(event) {
    event.preventDefault();
    
    const purchase = {
        product_id: parseInt(document.getElementById('purchaseProductId').value),
        quantity: parseInt(document.getElementById('purchaseQuantity').value),
        supplier: document.getElementById('purchaseSupplier').value,
        price: parseFloat(document.getElementById('purchasePrice').value)
    };
    
    try {
        const response = await fetch(`${API_URL}/purchases`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(purchase)
        });
        
        if (response.ok) {
            const data = await response.json();
            alert(`Purchase recorded successfully!\nTotal: ${data.total.toFixed(2)} ₸`);
            document.getElementById('purchaseForm').reset();
            loadPurchases();
            loadProducts();
        } else {
            const error = await response.text();
            alert('Error: ' + error);
        }
    } catch (error) {
        alert('Error recording purchase: ' + error.message);
    }
}