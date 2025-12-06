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
    // window.scrollTo(0, 0); // Optional: disabled for smoother browsing
}

function updateActiveNavLink() {
    // Get all nav buttons
    const navButtons = document.querySelectorAll('.nav-link');
    const sections = document.querySelectorAll('.section');

    // Find active section
    let activeSectionId = 'home';
    sections.forEach(section => {
        if (section.classList.contains('active')) {
            activeSectionId = section.id;
        }
    });

    // Update buttons
    navButtons.forEach(btn => {
        // Remove active class
        btn.classList.remove('active');

        // Add active class if it matches the function call
        // We match by checking the onclick attribute or ID mapping
        const onclick = btn.getAttribute('onclick');
        if (onclick && onclick.includes(`'${activeSectionId}'`)) {
            btn.classList.add('active');
        }
    });
}

// ... setupForms remains the same ...

async function loadProviders() {
    try {
        const res = await fetch(`${API_BASE}/providers?limit=20`);
        const data = await res.json();
        const providers = data.data || data;

        const container = document.getElementById('providersList');
        if (!providers || providers.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-search"></i>
                    <h3>No Providers Found</h3>
                    <p>We couldn't find any service providers at the moment.</p>
                </div>
            `;
            return;
        }

        container.innerHTML = `
            <div class="providers-grid">
                ${providers.map(p => `
                    <div class="card provider-card">
                        <div class="provider-header mb-2">
                             <h3>${p.user?.full_name || 'Service Provider'}</h3>
                             <span class="badge ${p.is_available ? 'badge-success' : 'badge-error'}">
                                ${p.is_available ? 'Available' : 'Busy'}
                             </span>
                        </div>
                        
                        <div class="provider-info">
                            <p class="text-muted mb-1"><i class="fas fa-map-marker-alt text-primary"></i> ${p.location || 'Remote'}</p>
                            <div class="provider-stats">
                                <span><i class="fas fa-star text-warning"></i> ${p.rating ? p.rating.toFixed(1) : 'New'}</span>
                                <span><i class="fas fa-briefcase text-info"></i> ${p.experience_years || 0}y exp</span>
                                <span><i class="fas fa-tag text-success"></i> $${p.hourly_rate || 0}/hr</span>
                            </div>
                        </div>
                        
                        ${p.bio ? `<p class="mt-3 text-sm">${p.bio}</p>` : ''}
                        
                        <button class="btn btn-primary btn-block mt-3 btn-small">Book Now</button>
                    </div>
                `).join('')}
            </div>
        `;
    } catch (error) {
        console.error(error);
        container.innerHTML = `<div class="alert alert-error">Failed to load providers</div>`;
    }
}

function renderStars(rating) {
    // Deprecated for new UI design, but kept for compatibility if needed
    return rating;
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

