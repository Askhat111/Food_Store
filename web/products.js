async function loadProducts() {
    try {
        const search = document.getElementById('searchProducts').value;
        const category = document.getElementById('categoryFilter').value;
        let url = `${API_URL}/products?`;
        if (search) url += `search=${search}&`;
        if (category) url += `category=${category}`;

        const response = await fetch(url);
        const products = await response.json();
        
        const tbody = document.getElementById('productsTableBody');
        tbody.innerHTML = '';
        
        products.forEach(product => {
            const tr = document.createElement('tr');
            
            //status badge for quantity
            let qtyBadge = '';
            const minStock = product.min_stock_threshold || 5;
            if (product.quantity < minStock) {
                qtyBadge = `<span class="status-badge status-low">${product.quantity}</span>`;
            } else {
                qtyBadge = `<span class="status-badge status-ok">${product.quantity}</span>`;
            }
            
            //status badge for expiry
            let expiryDisplay = product.expiry_date || '-';
            if (product.days_to_expiry !== undefined) {
                if (product.days_to_expiry < 0) {
                    expiryDisplay = `<span class="status-badge status-expiring">Expired</span>`;
                } else if (product.days_to_expiry <= 3) {
                    expiryDisplay = `<span class="status-badge status-expiring">${product.expiry_date} (${product.days_to_expiry}d)</span>`;
                } else if (product.days_to_expiry <= 7) {
                    expiryDisplay = `${product.expiry_date} (${product.days_to_expiry}d)`;
                }
            }
            
            let actions = '-';
            if (currentUser.role === 'owner') {
                actions = `
                    <button class="btn-edit" onclick="editProduct(${product.id})">Edit</button>
                    <button class="btn-delete" onclick="deleteProduct(${product.id})">Delete</button>
                `;
            }
            
            tr.innerHTML = `
                <td>${product.name}</td>
                <td>${product.category || '-'}</td>
                <td>${product.price.toFixed(2)} ₸</td>
                <td>${qtyBadge}</td>
                <td>${expiryDisplay}</td>
                <td>${product.min_stock_threshold || 5}</td>
                <td>${actions}</td>
            `;
            tbody.appendChild(tr);
        });
        
        updateProductDropdowns(products);
    } catch (error) {
        console.error('Error loading products:', error);
    }
}

function searchProducts() {
    loadProducts();
}

function filterProducts() {
    loadProducts();
}

function updateProductDropdowns(products) {
    const saleSelect = document.getElementById('saleProductId');
    const purchaseSelect = document.getElementById('purchaseProductId');
    
    saleSelect.innerHTML = '<option value="">Select Product</option>';
    purchaseSelect.innerHTML = '<option value="">Select Product</option>';
    
    products.forEach(product => {
        const option1 = document.createElement('option');
        option1.value = product.id;
        option1.textContent = `${product.name} (Stock: ${product.quantity})`;
        saleSelect.appendChild(option1);
        
        const option2 = document.createElement('option');
        option2.value = product.id;
        option2.textContent = product.name;
        purchaseSelect.appendChild(option2);
    });
}

function showAddProductForm() {
    document.getElementById('productForm').style.display = 'block';
    document.getElementById('formTitle').textContent = 'Add Product';
    document.getElementById('productFormElement').reset();
    document.getElementById('productId').value = '';
}

function hideProductForm() {
    document.getElementById('productForm').style.display = 'none';
    document.getElementById('productFormElement').reset();
}

async function saveProduct(event) {
    event.preventDefault();
    
    const id = document.getElementById('productId').value;
    const product = {
        name: document.getElementById('productName').value,
        category: document.getElementById('productCategory').value,
        price: parseFloat(document.getElementById('productPrice').value),
        quantity: parseInt(document.getElementById('productQuantity').value),
        expiry_date: document.getElementById('productExpiry').value,
        min_stock_threshold: parseInt(document.getElementById('productMinStock').value)
    };
    
    try {
        const url = id ? `${API_URL}/products/${id}` : `${API_URL}/products`;
        const method = id ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(product)
        });
        
        if (response.ok) {
            alert(id ? 'Product updated successfully!' : 'Product added successfully!');
            hideProductForm();
            loadProducts();
        } else {
            const error = await response.text();
            alert('Error: ' + error);
        }
    } catch (error) {
        alert('Error saving product: ' + error.message);
    }
}

async function editProduct(id) {
    try {
        const response = await fetch(`${API_URL}/products/${id}`);
        const product = await response.json();
        
        document.getElementById('productForm').style.display = 'block';
        document.getElementById('formTitle').textContent = 'Edit Product';
        document.getElementById('productId').value = product.id;
        document.getElementById('productName').value = product.name;
        document.getElementById('productCategory').value = product.category;
        document.getElementById('productPrice').value = product.price;
        document.getElementById('productQuantity').value = product.quantity;
        document.getElementById('productExpiry').value = product.expiry_date;
        document.getElementById('productMinStock').value = product.min_stock_threshold || 5;
        
        document.getElementById('productForm').scrollIntoView({ behavior: 'smooth' });
    } catch (error) {
        alert('Error loading product: ' + error.message);
    }
}

async function deleteProduct(id) {
    if (!confirm('Are you sure you want to delete this product?')) return;
    
    try {
        const response = await fetch(`${API_URL}/products/${id}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        
        if (response.ok) {
            alert('Product deleted successfully!');
            loadProducts();
        } else {
            const error = await response.text();
            alert('Error: ' + error);
        }
    } catch (error) {
        alert('Error deleting product: ' + error.message);
    }
}