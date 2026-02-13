document.addEventListener('DOMContentLoaded', () => {
    if (token && currentUser) {
        showApp();
    } else {
        showLogin();
    }
});

function showLogin() {
    document.getElementById('loginScreen').classList.add('active');
    document.getElementById('appScreen').classList.remove('active');
}

function showApp() {
    document.getElementById('loginScreen').classList.remove('active');
    document.getElementById('appScreen').classList.add('active');
    document.getElementById('userName').textContent = currentUser.name + ' (' + currentUser.role + ')';
    
    if (currentUser.role === 'staff') {
        document.getElementById('addProductBtn').style.display = 'none';
        document.querySelector('[onclick="showTab(\'purchases\')"]').style.display = 'none';
        document.querySelector('[onclick="showTab(\'reports\')"]').style.display = 'none';
    }
    
    loadProducts();
    loadAlerts();
    loadSales();
    loadPurchases();
    
    setInterval(loadAlerts, 30000);
}

document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        if (response.ok) {
            const data = await response.json();
            token = data.token;
            currentUser = data.user;
            localStorage.setItem('token', token);
            localStorage.setItem('user', JSON.stringify(currentUser));
            showApp();
        } else {
            alert('Invalid credentials');
        }
    } catch (error) {
        alert('Login error: ' + error.message);
    }
});

function logout() {
    token = null;
    currentUser = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    showLogin();
}

function showTab(tabName) {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(tc => tc.classList.remove('active'));
    
    //selected tab
    event.target.classList.add('active');
    document.getElementById(tabName + 'Tab').classList.add('active');
    
    if (tabName === 'products') loadProducts();
    if (tabName === 'sales') loadSales();
    if (tabName === 'purchases') loadPurchases();
    if (tabName === 'reports') loadReports();
}