let token = localStorage.getItem('token');
let currentUser = JSON.parse(localStorage.getItem('user') || 'null');
const API_URL = 'http://localhost:8080/api';