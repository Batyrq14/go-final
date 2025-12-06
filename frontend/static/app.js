const API_BASE = 'http://localhost:8080/api';
let authToken = localStorage.getItem('authToken');
let currentUser = null;

document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    setupForms();
    loadServices();
    loadProviders();

    // Restore section from URL hash
    const hash = window.location.hash.substring(1);
    if (hash && document.getElementById(hash)) {
        showSection(hash);
    } else {
        updateActiveNavLink();
    }
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

    // Role-specific UI changes
    const providersBtn = document.querySelector('button[onclick="showSection(\'providers\')"]');
    const servicesBtn = document.querySelector('button[onclick="showSection(\'services\')"]');

    if (currentUser.role === 'provider') {
        if (providersBtn) providersBtn.innerHTML = '<i class="fas fa-network-wired"></i> Competitors';
        if (servicesBtn) servicesBtn.innerHTML = '<i class="fas fa-briefcase"></i> Market Demand';
        document.getElementById('bookingsNavBtn').innerHTML = '<i class="fas fa-clipboard-list"></i> Requests';
    } else {
        // Client view defaults
        if (providersBtn) providersBtn.innerHTML = '<i class="fas fa-user-tie"></i> Find Pros';
        if (servicesBtn) servicesBtn.innerHTML = '<i class="fas fa-list"></i> Services';
        document.getElementById('bookingsNavBtn').innerHTML = '<i class="fas fa-calendar"></i> My Bookings';
    }
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

    updateActiveNavLink();

    history.pushState(null, null, `#${sectionId}`);

    // Load data when section is shown
    if (sectionId === 'bookings' && authToken) {
        loadBookings();
    }
    if (sectionId === 'profile' && authToken) {
        loadProfile();
    }

    // Scroll to top
    // window.scrollTo(0, 0); 
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
        btn.classList.remove('active');

        // Check onclick attribute to match section
        const onclick = btn.getAttribute('onclick');
        if (onclick && onclick.includes(`'${activeSectionId}'`)) {
            btn.classList.add('active');
        }
    });
}

function setupForms() {
    // Register Form
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('regEmail').value;
            const full_name = document.getElementById('regFullName').value;
            const phone = document.getElementById('regPhone').value;
            const password = document.getElementById('regPassword').value;
            const role = document.getElementById('regRole').value;

            try {
                const res = await fetch(`${API_BASE}/auth/register`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, full_name, phone, password, role })
                });
                const data = await res.json();

                if (res.ok) {
                    authToken = data.token;
                    localStorage.setItem('authToken', authToken);
                    currentUser = data.user;
                    updateUIForLoggedIn();
                    showSection('home');
                    document.getElementById('registerAlert').innerHTML =
                        '<div class="alert alert-success">Registration successful!</div>';
                } else {
                    document.getElementById('registerAlert').innerHTML =
                        `<div class="alert alert-error">${data.error || 'Registration failed'}</div>`;
                }
            } catch (err) {
                console.error(err);
                document.getElementById('registerAlert').innerHTML =
                    '<div class="alert alert-error">Network error</div>';
            }
        });
    }

    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('loginEmail').value;
            const password = document.getElementById('loginPassword').value;

            try {
                const res = await fetch(`${API_BASE}/auth/login`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });
                const data = await res.json();

                if (res.ok) {
                    authToken = data.token;
                    localStorage.setItem('authToken', authToken);
                    currentUser = data.user;
                    updateUIForLoggedIn();
                    showSection('home');
                } else {
                    document.getElementById('loginAlert').innerHTML =
                        `<div class="alert alert-error">${data.error || 'Login failed'}</div>`;
                }
            } catch (err) {
                console.error(err);
                document.getElementById('loginAlert').innerHTML =
                    '<div class="alert alert-error">Network error</div>';
            }
        });
    }
}

// Dropdown Toggle Logic
const userBtn = document.getElementById('userBtn');
const dropdownContent = document.querySelector('.dropdown-content');

if (userBtn) {
    userBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        dropdownContent.classList.toggle('show');
        userBtn.classList.toggle('active');
    });

    // Close on click outside
    document.addEventListener('click', (e) => {
        if (!userBtn.contains(e.target) && !dropdownContent.contains(e.target)) {
            dropdownContent.classList.remove('show');
            userBtn.classList.remove('active');
        }
    });
}


// ... existing loadServices ...

let allServices = [];

async function loadServices() {
    // Safety timeout to clear loading state
    const timeout = setTimeout(() => {
        const container = document.getElementById('servicesList');
        if (container && container.classList.contains('loading')) {
            container.innerHTML = '<p class="text-error">Request timed out.</p>';
        }
    }, 5000);

    try {
        const res = await fetch(`${API_BASE}/services`);
        clearTimeout(timeout); // Clear timeout on response
        const data = await res.json();

        const services = data.services || data.data || [];
        allServices = services;

        const container = document.getElementById('servicesList');
        if (!services || services.length === 0) {
            container.innerHTML = '<p>No services available.</p>';
            return;
        }

        container.innerHTML = `
            <div class="services-grid">
                ${services.map(s => `
                    <button class="card service-card" onclick="selectService('${s.id}')">
                        <i class="${s.icon_url || 'fas fa-concierge-bell'} service-icon"></i>
                        <h3>${s.name || s.title}</h3> 
                        <p>${s.description}</p>
                    </button>
                `).join('')}
            </div>
        `;
    } catch (error) {
        clearTimeout(timeout);
        console.error(error);
        const container = document.getElementById('servicesList');
        if (container) container.innerHTML = '<p class="text-error">Failed to load services</p>';
    }
}

// ... existing loadProfile ...

let allProviders = [];

async function loadProviders() {
    try {
        const res = await fetch(`${API_BASE}/providers?limit=20`);
        const data = await res.json();
        allProviders = data.providers || [];

        renderProviders(allProviders);
    } catch (error) {
        console.error(error);
        const container = document.getElementById('providersList');
        if (container) container.innerHTML = `<div class="alert alert-error">Failed to load providers</div>`;
    }
}

function renderProviders(providers) {
    const container = document.getElementById('providersList');

    // Apply Filters
    const searchValue = document.getElementById('providerSearch')?.value.toLowerCase();
    const ratingFilter = document.getElementById('providerFilter')?.value;

    const filtered = providers.filter(p => {
        const matchesSearch = !searchValue || (p.user.full_name.toLowerCase().includes(searchValue) || p.bio.toLowerCase().includes(searchValue));
        const matchesRating = !ratingFilter || (p.rating >= parseFloat(ratingFilter));
        return matchesSearch && matchesRating;
    });

    if (!filtered || filtered.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-search"></i>
                <h3>No Providers Found</h3>
            </div>
        `;
        return;
    }

    container.innerHTML = `
        <div class="providers-grid">
            ${filtered.map(p => `
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
                    
                    <button class="btn btn-primary btn-block mt-3 btn-small" onclick="initiateBooking('${p.provider_id}', '${p.user.full_name}')">Book Now</button>
                </div>
            `).join('')}
        </div>
    `;
}


document.getElementById('providerSearch')?.addEventListener('input', () => renderProviders(allProviders));
document.getElementById('providerFilter')?.addEventListener('change', () => renderProviders(allProviders));
document.querySelector('.search-btn')?.addEventListener('click', () => renderProviders(allProviders));


function selectService(serviceId) {
    showSection('providers');
    showSection('providers');
    window.selectedServiceId = serviceId;
    alert(`Selected service. Now choose a provider.`);
}

async function initiateBooking(providerId, providerName) {
    if (!authToken) {
        alert("Please login to book a provider.");
        showSection('login');
        return;
    }

    let serviceId = window.selectedServiceId;
    if (!serviceId && allServices.length > 0) {
        serviceId = allServices[0].id; // Fallback
    }

    if (!serviceId) {
        // Should ideally prompt user or create a "General" service if DB allows
        // But DB requires valid UUID referencing services table.
        // If no services exist, we are stuck.
        alert("No services available to book. Please contact admin.");
        return;
    }

    // Simple prompt for date (YYYY-MM-DD)
    const dateStr = prompt(`Booking with ${providerName}\nEnter date (YYYY-MM-DD):`, new Date().toISOString().split('T')[0]);
    if (!dateStr) return;

    // Default to 9 AM
    const scheduledTime = new Date(dateStr + 'T09:00:00Z').toISOString();

    try {
        const res = await fetch(`${API_BASE}/bookings`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({
                user_id: currentUser.id, // Gateway injects this usually, but safe to send
                provider_id: providerId,
                service_id: serviceId,
                scheduled_time: scheduledTime
            })
        });

        if (res.ok) {
            alert("Booking sent successfully! You can view it in 'My Bookings'.");
            showSection('bookings');
        } else {
            const err = await res.json();
            alert("Failed to book: " + (err.error || "Unknown error"));
        }
    } catch (e) {
        console.error(e);
        alert("Network error while booking.");
    }
}


// ... existing Smooth scroll ...

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

        const bookings = data.bookings || [];

        const container = document.getElementById('bookingsList');
        if (!bookings || bookings.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-inbox"></i>
                    <h3>No Bookings Yet</h3>
                    <p>You haven't made or received any bookings yet.</p>
                </div>
            `;
            return;
        }

        const isProvider = currentUser.role === 'provider';

        container.innerHTML = `
            <div class="bookings-grid">
                ${bookings.map(b => `
                    <div class="card booking-card">
                        <div class="booking-header mb-2">
                             <h3><i class="fas fa-calendar-alt"></i> ${new Date(b.scheduled_time).toLocaleDateString()}</h3>
                             <span class="badge badge-${getStatusColor(b.status)}">
                                <i class="fas ${getStatusIcon(b.status)}"></i>
                                ${b.status}
                             </span>
                        </div>
                        <p><i class="fas fa-clock"></i> <strong>Time:</strong> ${new Date(b.scheduled_time).toLocaleTimeString()}</p>
                        <p><i class="fas fa-tools"></i> <strong>Service:</strong> ${b.service_title || b.service_id}</p>
                        <p><i class="fas ${isProvider ? 'fa-user' : 'fa-user-tie'}"></i> <strong>${isProvider ? 'Client' : 'Provider'}:</strong> ${b.other_party_name || (isProvider ? b.client_id : b.provider_id)}</p>
                        ${b.status === 'pending' && isProvider ? `
                            <div class="booking-actions mt-3">
                                <button class="btn btn-success btn-small" onclick="updateBookingStatus('${b.id}', 'confirmed')"><i class="fas fa-check"></i> Accept</button>
                                <button class="btn btn-error btn-small" onclick="updateBookingStatus('${b.id}', 'rejected')"><i class="fas fa-times"></i> Reject</button>
                            </div>
                        ` : ''}
                        ${b.status === 'confirmed' && isProvider ? `
                             <div class="booking-actions mt-3">
                                <button class="btn btn-primary btn-small" onclick="updateBookingStatus('${b.id}', 'completed')"><i class="fas fa-check-double"></i> Mark Completed</button>
                            </div>
                        ` : ''}
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

async function updateBookingStatus(bookingId, status) {
    if (!confirm(`Are you sure you want to ${status} this booking?`)) return;

    try {
        const res = await fetch(`${API_BASE}/bookings/${bookingId}/status`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({ status })
        });

        if (res.ok) {
            loadBookings(); // Reload list
        } else {
            const err = await res.json();
            alert("Failed to update status: " + (err.error || "Unknown error"));
        }
    } catch (e) {
        console.error(e);
        alert("Network error.");
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
            <div class="profile-header-card card">
                <div class="profile-avatar">
                    <span class="avatar-initials">${user.full_name ? user.full_name.charAt(0).toUpperCase() : user.email.charAt(0).toUpperCase()}</span>
                </div>
                <div class="profile-info-main">
                    <h2>${user.full_name || 'User'}</h2>
                    <p class="text-muted"><i class="fas fa-envelope"></i> ${user.email}</p>
                     <span class="badge badge-primary mt-2">
                        <i class="fas ${user.role === 'provider' ? 'fa-user-tie' : 'fa-user'}"></i>
                        ${user.role === 'provider' ? 'Service Provider' : 'Client'}
                    </span>
                </div>
            </div>
            
            <div class="profile-details-grid mt-4">
                 <div class="card detail-card">
                    <div class="detail-icon"><i class="fas fa-phone"></i></div>
                    <div class="detail-content">
                        <label>Phone</label>
                        <p>${user.phone || 'Not provided'}</p>
                    </div>
                </div>
                 <div class="card detail-card">
                    <div class="detail-icon"><i class="fas fa-calendar"></i></div>
                    <div class="detail-content">
                        <label>Joined</label>
                        <p>${createdDate}</p>
                    </div>
                </div>
                  <div class="card detail-card">
                    <div class="detail-icon"><i class="fas fa-shield-alt"></i></div>
                    <div class="detail-content">
                        <label>Account ID</label>
                        <p class="text-xs text-muted">${user.id}</p>
                    </div>
                </div>
            </div>

            <div class="mt-4 text-center">
                 <button class="btn btn-secondary" onclick="logout()">
                    <i class="fas fa-sign-out-alt"></i> Logout
                 </button>
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

