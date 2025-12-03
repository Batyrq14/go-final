const API_BASE = 'http://localhost:8080/api';
let authToken = localStorage.getItem('authToken');
let currentUser = null;

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    setupForms();
    loadServices();
    loadProviders();
    updateActiveNavLink();
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
                if (data.id || data.email) {
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
    document.getElementById('userDisplayName').textContent = currentUser.email || 'User';
    document.getElementById('userDropdown').classList.remove('hidden');
    document.getElementById('registerNavBtn').classList.add('hidden');
    document.getElementById('loginNavBtn').classList.add('hidden');
    document.getElementById('bookingsNavBtn').classList.remove('hidden');
    document.getElementById('profileNavBtn').classList.remove('hidden');
}

function logout() {
    authToken = null;
    currentUser = null;
    localStorage.removeItem('authToken');
    document.getElementById('userDropdown').classList.add('hidden');
    document.getElementById('registerNavBtn').classList.remove('hidden');
    document.getElementById('loginNavBtn').classList.remove('hidden');
    document.getElementById('bookingsNavBtn').classList.add('hidden');
    document.getElementById('profileNavBtn').classList.add('hidden');
    showSection('home');
}

function showSection(sectionId) {
    // Hide all sections
    document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
    document.getElementById(sectionId).classList.add('active');

    // Update nav links
    updateActiveNavLink();

    // Load data when section is shown
    if (sectionId === 'bookings' && authToken) {
        loadBookings();
    }
    if (sectionId === 'profile' && authToken) {
        loadProfile();
    }

    // Scroll to top
    window.scrollTo(0, 0);
}

function updateActiveNavLink() {
    // This would need to be updated based on current section
    // For simplicity, we'll handle it with active classes in buttons
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
            console.error(error);
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
            console.error(error);
        }
    });
}

function showAlert(containerId, message, type) {
    const container = document.getElementById(containerId);
    const icon = type === 'success' ? 'fa-check-circle' : type === 'error' ? 'fa-exclamation-circle' : 'fa-info-circle';
    container.innerHTML = `<div class="alert alert-${type}"><i class="fas ${icon}"></i>${message}</div>`;
    setTimeout(() => container.innerHTML = '', 5000);
}

async function loadServices() {
    try {
        const res = await fetch(`${API_BASE}/services`);
        const data = await res.json();
        const services = data.data || data;

        const container = document.getElementById('servicesList');
        if (!services || services.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-inbox"></i>
                    <h3>No Services Available</h3>
                    <p>Services will be available soon. Check back later!</p>
                </div>
            `;
            return;
        }

        container.innerHTML = `
            <div class="services-grid">
                ${services.map(s => `
                    <div class="card">
                        <h3><i class="fas fa-concierge-bell"></i> ${s.name}</h3>
                        <p><strong>Description:</strong></p>
                        <p>${s.description || 'No description available'}</p>
                    </div>
                `).join('')}
            </div>
        `;
    } catch (error) {
        document.getElementById('servicesList').innerHTML = `
            <div class="alert alert-error">
                <i class="fas fa-exclamation-triangle"></i>
                Failed to load services. Make sure the API is running on http://localhost:8080
            </div>
        `;
        console.error(error);
    }
}

async function loadProviders() {
    try {
        const res = await fetch(`${API_BASE}/providers?limit=20`);
        const data = await res.json();
        const providers = data.data || data;

        const container = document.getElementById('providersList');
        if (!providers || providers.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-inbox"></i>
                    <h3>No Providers Available</h3>
                    <p>No service providers found. Check back soon!</p>
                </div>
            `;
            return;
        }

        container.innerHTML = `
            <div class="providers-grid">
                ${providers.map(p => `
                    <div class="card">
                        <h3><i class="fas fa-user-tie"></i> ${p.user?.full_name || 'Service Provider'}</h3>
                        <p><i class="fas fa-map-marker-alt"></i> <strong>Location:</strong> ${p.location || 'Not specified'}</p>
                        <p><i class="fas fa-dollar-sign"></i> <strong>Rate:</strong> $${p.hourly_rate || '0'}/hour</p>
                        <p><i class="fas fa-briefcase"></i> <strong>Experience:</strong> ${p.experience_years || '0'} years</p>
                        <p>
                            <strong>Rating:</strong>
                            <span class="stars">${renderStars(p.rating || 0)}</span>
                            ${p.rating ? `<span class="rating-text">${p.rating.toFixed(1)}</span>` : '<span class="text-primary">No ratings yet</span>'}
                        </p>
                        <p>
                            <strong>Availability:</strong>
                            <span class="badge ${p.is_available ? 'badge-success' : 'badge-error'}">
                                <i class="fas ${p.is_available ? 'fa-check-circle' : 'fa-times-circle'}"></i>
                                ${p.is_available ? 'Available' : 'Unavailable'}
                            </span>
                        </p>
                        ${p.bio ? `<p><i class="fas fa-align-left"></i> <strong>Bio:</strong> ${p.bio}</p>` : ''}
                    </div>
                `).join('')}
            </div>
        `;
    } catch (error) {
        document.getElementById('providersList').innerHTML = `
            <div class="alert alert-error">
                <i class="fas fa-exclamation-triangle"></i>
                Failed to load providers. Make sure the API is running.
            </div>
        `;
        console.error(error);
    }
}

function renderStars(rating) {
    const fullStars = Math.round(rating);
    return '⭐'.repeat(Math.min(fullStars, 5)) || 'No rating';
}

async function loadBookings() {
    if (!authToken) {
        document.getElementById('bookingsList').innerHTML = `
            <div class="alert alert-info">
                <i class="fas fa-lock"></i>
                Please login to view your bookings
            </div>
        `;
        return;
    }

    try {
        const res = await fetch(`${API_BASE}/bookings`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        const data = await res.json();
        const bookings = data.data || data;

        const container = document.getElementById('bookingsList');
        if (!bookings || bookings.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-inbox"></i>
                    <h3>No Bookings Yet</h3>
                    <p>You haven't made any bookings yet. Browse providers to get started!</p>
                </div>
            `;
            return;
        }

        container.innerHTML = `
            <div class="bookings-grid">
                ${bookings.map(b => `
                    <div class="card">
                        <h3><i class="fas fa-calendar-alt"></i> Booking #${b.id.substring(0, 8)}</h3>
                        <p><i class="fas fa-tools"></i> <strong>Service:</strong> ${b.service?.name || 'N/A'}</p>
                        <p><i class="fas fa-user"></i> <strong>Provider:</strong> ${b.provider?.user?.full_name || 'N/A'}</p>
                        <p><i class="fas fa-clock"></i> <strong>Date:</strong> ${new Date(b.scheduled_date).toLocaleString()}</p>
                        <p><i class="fas fa-hourglass-half"></i> <strong>Duration:</strong> ${b.duration_hours} hours</p>
                        <p><i class="fas fa-dollar-sign"></i> <strong>Total:</strong> $${b.total_price || '0'}</p>
                        <p>
                            <strong>Status:</strong>
                            <span class="badge badge-${getStatusColor(b.status)}">
                                <i class="fas ${getStatusIcon(b.status)}"></i>
                                ${b.status}
                            </span>
                        </p>
                        ${b.notes ? `<p><i class="fas fa-sticky-note"></i> <strong>Notes:</strong> ${b.notes}</p>` : ''}
                    </div>
                `).join('')}
            </div>
        `;
    } catch (error) {
        document.getElementById('bookingsList').innerHTML = `
            <div class="alert alert-error">
                <i class="fas fa-exclamation-triangle"></i>
                Failed to load bookings
            </div>
        `;
        console.error(error);
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
    return colors[status] || 'info';
}

function getStatusIcon(status) {
    const icons = {
        'pending': 'fa-clock',
        'accepted': 'fa-check-circle',
        'rejected': 'fa-times-circle',
        'completed': 'fa-check-double',
        'cancelled': 'fa-ban'
    };
    return icons[status] || 'fa-question-circle';
}

async function loadProfile() {
    if (!authToken) {
        document.getElementById('profileContent').innerHTML = `
            <div class="alert alert-info">
                <i class="fas fa-lock"></i>
                Please login to view your profile
            </div>
        `;
        return;
    }

    try {
        const res = await fetch(`${API_BASE}/auth/me`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        const user = await res.json();

        const createdDate = new Date(user.created_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });

        document.getElementById('profileContent').innerHTML = `
            <div class="card">
                <h3><i class="fas fa-user-circle"></i> Profile Information</h3>
                <p><i class="fas fa-envelope"></i> <strong>Email:</strong> ${user.email}</p>
                <p><i class="fas fa-user"></i> <strong>Full Name:</strong> ${user.full_name || 'Not provided'}</p>
                <p><i class="fas fa-phone"></i> <strong>Phone:</strong> ${user.phone || 'Not provided'}</p>
                <p>
                    <strong>Role:</strong>
                    <span class="badge badge-primary">
                        <i class="fas ${user.role === 'provider' ? 'fa-user-tie' : 'fa-user'}"></i>
                        ${user.role === 'provider' ? 'Service Provider' : 'Client'}
                    </span>
                </p>
                <p><i class="fas fa-calendar"></i> <strong>Member Since:</strong> ${createdDate}</p>
            </div>
        `;
    } catch (error) {
        document.getElementById('profileContent').innerHTML = `
            <div class="alert alert-error">
                <i class="fas fa-exclamation-triangle"></i>
                Failed to load profile
            </div>
        `;
        console.error(error);
    }
}

// Smooth scroll for buttons
document.querySelectorAll('a[href^="#"], button').forEach(btn => {
    btn.addEventListener('click', function (e) {
        if (this.getAttribute('href')?.startsWith('#')) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({ behavior: 'smooth' });
            }
        }
    });
});

