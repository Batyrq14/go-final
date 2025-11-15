const API_BASE = 'http://localhost:8080/api';
let authToken = localStorage.getItem('authToken');
let currentUser = null;

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    setupForms();
    loadServices();
    loadProviders();
});

function checkAuth() {
    if (authToken) {
        fetch(`${API_BASE}/auth/me`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        })
        .then(res => res.json())
        .then(data => {
            if (data.id) {
                currentUser = data;
                updateUIForLoggedIn();
            } else {
                logout();
            }
        })
        .catch(() => logout());
    }
}

function updateUIForLoggedIn() {
    document.getElementById('userEmail').textContent = currentUser.email;
    document.getElementById('userRole').textContent = currentUser.role;
    document.getElementById('userInfo').classList.add('active');
    document.getElementById('registerBtn').classList.add('hidden');
    document.getElementById('loginBtn').classList.add('hidden');
    document.getElementById('bookingsBtn').classList.remove('hidden');
    document.getElementById('profileBtn').classList.remove('hidden');
}

function logout() {
    authToken = null;
    currentUser = null;
    localStorage.removeItem('authToken');
    document.getElementById('userInfo').classList.remove('active');
    document.getElementById('registerBtn').classList.remove('hidden');
    document.getElementById('loginBtn').classList.remove('hidden');
    document.getElementById('bookingsBtn').classList.add('hidden');
    document.getElementById('profileBtn').classList.add('hidden');
    showSection('home');
}

function showSection(sectionId) {
    document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
    document.getElementById(sectionId).classList.add('active');
    
    document.querySelectorAll('.nav button').forEach(b => b.classList.remove('active'));
    event?.target?.classList.add('active');

    if (sectionId === 'bookings' && authToken) {
        loadBookings();
    }
    if (sectionId === 'profile' && authToken) {
        loadProfile();
    }
}

function setupForms() {
    document.getElementById('registerForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const data = {
            email: document.getElementById('regEmail').value,
            password: document.getElementById('regPassword').value,
            full_name: document.getElementById('regFullName').value,
            phone: document.getElementById('regPhone').value,
            role: document.getElementById('regRole').value
        };

        try {
            const res = await fetch(`${API_BASE}/auth/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            const result = await res.json();
            
            if (res.ok) {
                showAlert('registerAlert', 'Registration successful! Please login.', 'success');
                document.getElementById('registerForm').reset();
                setTimeout(() => showSection('login'), 2000);
            } else {
                showAlert('registerAlert', result.error || 'Registration failed', 'error');
            }
        } catch (error) {
            showAlert('registerAlert', 'Network error. Make sure the API is running.', 'error');
        }
    });

    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const data = {
            email: document.getElementById('loginEmail').value,
            password: document.getElementById('loginPassword').value
        };

        try {
            const res = await fetch(`${API_BASE}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            const result = await res.json();
            
            if (res.ok && result.token) {
                authToken = result.token;
                currentUser = result.user;
                localStorage.setItem('authToken', authToken);
                showAlert('loginAlert', 'Login successful!', 'success');
                document.getElementById('loginForm').reset();
                updateUIForLoggedIn();
                setTimeout(() => showSection('home'), 1000);
            } else {
                showAlert('loginAlert', result.error || 'Login failed', 'error');
            }
        } catch (error) {
            showAlert('loginAlert', 'Network error. Make sure the API is running.', 'error');
        }
    });
}

function showAlert(containerId, message, type) {
    const container = document.getElementById(containerId);
    container.innerHTML = `<div class="alert alert-${type}">${message}</div>`;
    setTimeout(() => container.innerHTML = '', 5000);
}

async function loadServices() {
    try {
        const res = await fetch(`${API_BASE}/services`);
        const services = await res.json();
        
        const container = document.getElementById('servicesList');
        if (services.length === 0) {
            container.innerHTML = '<p>No services available yet.</p>';
            return;
        }
        
        container.innerHTML = services.map(s => `
            <div class="card">
                <h3>${s.name}</h3>
                <p>${s.description || 'No description'}</p>
            </div>
        `).join('');
    } catch (error) {
        document.getElementById('servicesList').innerHTML = 
            '<p class="alert alert-error">Failed to load services. Make sure the API is running.</p>';
    }
}

async function loadProviders() {
    try {
        const res = await fetch(`${API_BASE}/providers?limit=10`);
        const result = await res.json();
        const providers = result.data || result;
        
        const container = document.getElementById('providersList');
        if (!providers || providers.length === 0) {
            container.innerHTML = '<p>No providers available yet.</p>';
            return;
        }
        
        container.innerHTML = providers.map(p => `
            <div class="card">
                <h3>${p.user?.full_name || 'Provider'}</h3>
                <p><strong>Location:</strong> ${p.location || 'Not specified'}</p>
                <p><strong>Rate:</strong> $${p.hourly_rate}/hour</p>
                <p><strong>Experience:</strong> ${p.experience_years} years</p>
                <p><strong>Rating:</strong> ${p.rating ? '⭐'.repeat(Math.round(p.rating)) : 'No ratings yet'}</p>
                <p><strong>Available:</strong> ${p.is_available ? '✅ Yes' : '❌ No'}</p>
                ${p.bio ? `<p>${p.bio}</p>` : ''}
            </div>
        `).join('');
    } catch (error) {
        document.getElementById('providersList').innerHTML = 
            '<p class="alert alert-error">Failed to load providers. Make sure the API is running.</p>';
    }
}

async function loadBookings() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API_BASE}/bookings`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        const bookings = await res.json();
        
        const container = document.getElementById('bookingsList');
        if (!bookings || bookings.length === 0) {
            container.innerHTML = '<p>You have no bookings yet.</p>';
            return;
        }
        
        container.innerHTML = bookings.map(b => `
            <div class="card">
                <h3>Booking #${b.id.substring(0, 8)}</h3>
                <p><strong>Service:</strong> ${b.service?.name || 'N/A'}</p>
                <p><strong>Provider:</strong> ${b.provider?.user?.full_name || 'N/A'}</p>
                <p><strong>Date:</strong> ${new Date(b.scheduled_date).toLocaleString()}</p>
                <p><strong>Duration:</strong> ${b.duration_hours} hours</p>
                <p><strong>Total:</strong> $${b.total_price}</p>
                <p><strong>Status:</strong> <span class="badge badge-${getStatusColor(b.status)}">${b.status}</span></p>
            </div>
        `).join('');
    } catch (error) {
        document.getElementById('bookingsList').innerHTML = 
            '<p class="alert alert-error">Failed to load bookings.</p>';
    }
}

function getStatusColor(status) {
    const colors = {
        'pending': 'warning',
        'accepted': 'success',
        'rejected': 'error',
        'completed': 'success',
        'cancelled': 'error'
    };
    return colors[status] || 'primary';
}

async function loadProfile() {
    if (!authToken) return;
    
    try {
        const res = await fetch(`${API_BASE}/auth/me`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        const user = await res.json();
        
        document.getElementById('profileContent').innerHTML = `
            <div class="card">
                <h3>User Information</h3>
                <p><strong>Email:</strong> ${user.email}</p>
                <p><strong>Name:</strong> ${user.full_name}</p>
                <p><strong>Phone:</strong> ${user.phone || 'Not provided'}</p>
                <p><strong>Role:</strong> <span class="badge badge-primary">${user.role}</span></p>
                <p><strong>Member since:</strong> ${new Date(user.created_at).toLocaleDateString()}</p>
            </div>
        `;
    } catch (error) {
        document.getElementById('profileContent').innerHTML = 
            '<p class="alert alert-error">Failed to load profile.</p>';
    }
}

