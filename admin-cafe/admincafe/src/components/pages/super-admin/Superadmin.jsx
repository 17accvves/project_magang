import React, { useState } from 'react';
import './SuperAdmin.css';

// ==========================================
// 1. KOMPONEN KECIL (UI ELEMENTS)
// ==========================================

// Kartu Statistik (Kotak-kotak angka di Dashboard)
const StatCard = ({ title, value, icon, colorClass }) => (
    <div className={`sa-stat-card ${colorClass}`}>
        <div className="stat-info">
            <p>{title}</p>
            <h2>{value}</h2>
        </div>
        <div className="stat-icon">{icon}</div>
    </div>
);

// Kartu Cafe (Kotak-kotak di halaman Cafe)
const CafeCard = ({ name, rank, reviews, logo, rating }) => (
  <div className="cafe-card">
    <div className="cafe-logo-area">
      <img src={logo} alt={name} className="cafe-logo-img" />
      <button className="img-nav left">‹</button>
      <button className="img-nav right">›</button>
    </div>
    <div className="cafe-info">
      <div className="cafe-header-row">
        <h3>{name}</h3>
        <button className="btn-icon-edit">✏</button>
      </div>
      <p className="cafe-rank">#{rank} Bulan September</p>
      <div className="cafe-rating">
        <span className="stars">{'★'.repeat(Math.floor(rating))}</span>
        <span className="review-count">({reviews})</span>
      </div>
      <button className="btn-see-cafe">Lihat Cafe</button>
    </div>
  </div>
);

const DashboardContent = () => {
    const stats = [
        { title: 'User', value: '01171', icon: '👤', colorClass: 'bg-purple' },
        { title: 'Kafe Terdaftar', value: '00317', icon: '☕', colorClass: 'bg-orange' },
        { title: 'Event Organizer', value: '00317', icon: '📅', colorClass: 'bg-blue' },
        { title: 'Admin Aktif', value: '004', icon: '👮', colorClass: 'bg-green' },
    ];

    return (
        <div className="sa-content-padding">
            <h1 className="page-title">Dashboard</h1>
            
            {/* Statistik */}
            <div className="sa-grid-stats">
                {stats.map((s, i) => <StatCard key={i} {...s} />)}
            </div>

            {/* Grafik & Kalender */}
            <div className="sa-grid-main">
                <div className="sa-card chart-section">
                    <div className="card-header">
                        <h3>Pengunjung</h3>
                        <button className="btn-blue">Lihat</button>
                    </div>
                    <div className="chart-placeholder">
                        <div className="chart-lines">📈 [Grafik Visual Placeholder]</div>
                        <div className="chart-legend">
                            <span>🟣 Loyal</span> <span>🔴 Baru</span> <span>🟢 Premium</span>
                        </div>
                    </div>
                </div>

                <div className="sa-card calendar-section">
                    <h3>Sep, 19 Friday</h3>
                    <div className="calendar-placeholder">
                       <div className="cal-grid">
                           <span>M</span><span>T</span><span>W</span><span>T</span><span>F</span><span className="red">S</span><span className="red">S</span>
                           <span>18</span><span className="today">19</span><span>20</span><span>21</span><span>22</span><span className="red">23</span><span className="red">24</span>
                       </div>
                    </div>
                </div>
            </div>

            {/* Tabel */}
            <div className="sa-card table-section">
                <div className="card-header">
                    <h3>CAFE TERFAVORIT</h3>
                    <select><option>September</option></select>
                </div>
                <table className="sa-table">
                    <thead>
                        <tr>
                            <th>Nama Cafe</th>
                            <th>Lokasi</th>
                            <th>Total</th>
                            <th>Menu Favorit</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>☕ YOTTA.id</td>
                            <td>Jl. Pettarani</td>
                            <td>120</td>
                            <td>Extra Joss Susu</td>
                            <td><span className="badge active">Aktif</span></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    );
};

// --- KONTEN CAFE ---
const CafeContent = () => {
  const cafes = [
    { name: 'YOTTA.id', rank: 1, rating: 5, reviews: 131, logo: 'https://via.placeholder.com/150/2ecc71/ffffff?text=YOTTA' },
    { name: 'Fore Coffee', rank: 2, rating: 4, reviews: 120, logo: 'https://via.placeholder.com/150/ffffff/2ecc71?text=fore' },
    { name: 'Starbucks', rank: 3, rating: 4.5, reviews: 117, logo: 'https://via.placeholder.com/150/006241/ffffff?text=Starbucks' },
  ];

  return (
    <div className="sa-content-padding cafe-theme"> 
      <div className="page-header-row">
        <h1 className="page-title">Cafe</h1>
        <button className="btn-add-cafe">
            <span className="plus">+</span> Daftarkan Cafe
        </button>
      </div>

      <div className="promo-banner">
        <div className="banner-nav left">‹</div>
        <div className="banner-content">
            <span className="banner-date">September 25-30</span>
            <h2 className="banner-title">Nikmati keseruan LIVE MUSIC <br /> dengan artis di daerah kamu!</h2>
            <p className="banner-subtitle">• CARI Spot Super ADMIN</p>
            <button className="btn-promo-action">Atur Promo</button>
        </div>
        <div className="banner-nav right">›</div>
      </div>

      <div className="cafe-grid">
        {cafes.map((cafe, index) => (
            <CafeCard key={index} {...cafe} />
        ))}
      </div>
    </div>
  );
};

// --- KONTEN EVENT ---
const EventContent = () => {
  // Data Dummy untuk Event Berjalan
  const runningEvents = [
    { id: 1, cafe: 'Yotta.id', desc: 'band : Avenged 7x', eventName: 'Kelayapan', organizer: 'SMA Teykom', date: '19-09-2025', time: '17.00 WITA-Selesai', type: 'music' },
    { id: 2, cafe: 'Ngopisantuy', desc: 'band : UNM Band', eventName: 'Jam Malam', organizer: 'SMA Teykom', date: '19-09-2025', time: '14.00 WITA-Selesai', type: 'coffee' },
    { id: 3, cafe: 'Teh Indonesia', desc: 'band : free to all', eventName: 'spot.in', organizer: 'SMA Teykom', date: '19-09-2025', time: '22.00 WITA-Selesai', type: 'drink' },
    { id: 4, cafe: 'Yotta.id', desc: 'band : Avenged 7x', eventName: 'Kelayapan', organizer: 'SMA Teykom', date: '19-09-2025', time: '17.00 WITA-Selesai', type: 'music' },
    { id: 5, cafe: 'Ngopisantuy', desc: 'band : UNM Band', eventName: 'Jam Malam', organizer: 'SMA Teykom', date: '19-09-2025', time: '14.00 WITA-Selesai', type: 'coffee' },
    { id: 6, cafe: 'Teh Indonesia', desc: 'band : free to all', eventName: 'spot.in', organizer: 'SMA Teykom', date: '19-09-2025', time: '22.00 WITA-Selesai', type: 'drink' },
  ];

  return (
    <div className="sa-content-padding cafe-theme">
      <h1 className="page-title">Event</h1>

      <div className="event-layout-grid">
        
        {/* --- KOLOM KIRI --- */}
        <div className="event-left-col">
          
          {/* Widget EO Aktif */}
          <div className="sa-card eo-active-card">
            <div className="eo-info">
              <p className="label">EO Aktif</p>
              <h2 className="value">00074</h2>
            </div>
            <button className="btn-view-eo">
              📄 Lihat EO
            </button>
          </div>

          {/* Widget Event Berjalan */}
          <div className="sa-card running-event-card">
            <div className="card-header-center">
                <h3>Event Berjalan</h3>
                <div className="scroll-arrow up">︿</div>
            </div>

            <div className="event-table-header">
                <div className="col-header">
                    <span className="header-icon">🎵</span> Live Music
                </div>
                <div className="col-header">
                    <span className="header-icon">🍔</span> Bazar
                </div>
                <div className="col-header">
                    <span className="header-icon">🕒</span> Waktu
                </div>
            </div>

            <div className="event-list-container">
                {runningEvents.map((item, index) => (
                    <div className="event-row" key={index}>
                        {/* Kolom 1: Cafe / Musik */}
                        <div className="event-col col-1">
                            <span className="row-num">{item.id}</span>
                            <div className="col-icon type-icon">
                                {item.type === 'music' ? '🥤' : item.type === 'coffee' ? '☕' : '🥤'}
                            </div>
                            <div className="col-text">
                                <strong>{item.cafe}</strong>
                                <span>{item.desc}</span>
                            </div>
                        </div>

                        {/* Kolom 2: Event / Bazar */}
                        <div className="event-col col-2">
                            <div className="col-icon">🛵</div>
                            <div className="col-text">
                                <strong>{item.eventName}</strong>
                                <span>penyelenggara : {item.organizer}</span>
                            </div>
                        </div>

                        {/* Kolom 3: Waktu */}
                        <div className="event-col col-3">
                            <div className="col-text align-right">
                                <strong>{item.date}</strong>
                                <span>Pukul : {item.time}</span>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
            
            <div className="card-footer-center">
                <div className="scroll-arrow down">﹀</div>
            </div>
          </div>

        </div>

        {/* --- KOLOM KANAN --- */}
        <div className="event-right-col">
          
          {/* Widget Registrasi */}
          <div className="sa-card registration-card">
            <h3>Registrasi</h3>
            <div className="reg-buttons">
                <button className="btn-reg blue">Verifikasi Akun EO</button>
                <button className="btn-reg blue">Kelola Akun EO</button>
            </div>
          </div>

          {/* Widget Manage Event */}
          <div className="sa-card manage-event-card">
            <h3>Manage Event</h3>
            <div className="manage-item">
                <span>Edit</span>
                <button className="btn-icon">✏</button>
            </div>
            <div className="manage-item">
                <span>Lihat Riwayat</span>
                <button className="btn-icon">↺</button>
            </div>
          </div>

          {/* Widget Highlight Event */}
          <div className="sa-card highlight-card">
            <div className="highlight-img">
                {/* Ganti src dengan gambar event asli */}
                <img src="https://via.placeholder.com/300x150" alt="Event" />
                <div className="tag-yellow">Yotta.id</div>
            </div>
            <div className="highlight-info">
                <div className="date-box">
                    <span className="month">AUG</span>
                    <span className="date">27</span>
                </div>
                <div className="detail-box">
                    <h4>Bazar SMAN 1 Kilo</h4>
                    <p>Outlet Pettarani, Makassar</p>
                    <p className="time">14:20 WITA</p>
                    <div className="attendees">
                        👥 34 • ⭐⭐⭐⭐⭐
                    </div>
                </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
};

// --- KONTEN LOG AKTIVITAS ---

// Komponen Kartu Log (Reusable)
const ActivityCard = ({ title, date, status, subject, claim, variant }) => {
  // variant bisa 'dark' atau 'light'
  return (
    <div className={`activity-card ${variant}`}>
      <div className="act-content">
        <h3 className="act-title">{title}</h3>
        
        <div className="act-meta">
            <span className="act-date">{date}</span>
            <div className="act-status">
                Operasional : <span className="status-pill">{status}</span>
            </div>
        </div>

        <div className="act-details">
            <p>Perihal : {subject}</p>
            <p>{claim}</p>
        </div>
      </div>

      <div className="act-actions">
        <button className="btn-act confirm">
            <span className="icon-check">✔</span> Konfirmasi
        </button>
        <button className="btn-act history">
            <span className="icon-dots">•••</span> Riwayat
        </button>
        <button className="btn-act chat">
            <span className="icon-chat">💬</span> Buka Chat
        </button>
      </div>
    </div>
  );
};

const LogAktivitasContent = () => {
  // Data Dummy sesuai gambar
  const logs = [
    { 
      title: 'Kafe xxx mengirimkan perubahan Diskon', 
      date: '7 Sep 2025', status: 'Aktif', 
      subject: 'Diskon 20% sd 20rb', claim: 'Klaim : Anggota Premium', 
      variant: 'dark' 
    },
    { 
      title: 'Kafe xxx mengirimkan perubahan Lokasi', 
      date: '7 Sep 2025', status: 'Aktif', 
      subject: 'Perubahan Lokasi', claim: 'Klaim : SUPER ADMIN', 
      variant: 'light' 
    },
    { 
      title: 'Kafe xxx melakukan perubahan foto profil', 
      date: '7 Sep 2025', status: 'Aktif', 
      subject: 'Perubahan Akun', claim: '', 
      variant: 'light' 
    },
    { 
      title: 'USER Delon berlangganan PREMIUM', 
      date: '7 Sep 2025', status: 'Aktif', 
      subject: 'Berlangganan', claim: 'Klaim : Anggota Premium', 
      variant: 'dark' 
    },
    { 
      title: 'Kafe xxx mengirimkan perubahan Diskon', 
      date: '7 Sep 2025', status: 'Aktif', 
      subject: 'Diskon 20% sd 20rb', claim: 'Klaim : Anggota Premium', 
      variant: 'dark' // Ini yang ada border biru di gambar (selected)
    },
    { 
      title: 'Kafe xxx mengirimkan perubahan Lokasi', 
      date: '7 Sep 2025', status: 'Aktif', 
      subject: 'Perubahan Lokasi', claim: 'Datang : 75', 
      variant: 'light' 
    },
  ];

  return (
    <div className="sa-content-padding cafe-theme">
      <div className="log-header">
        <h1 className="page-title">Log Aktivitas</h1>
        <div className="log-header-actions">
            <button className="btn-header-green delete">🗑 Hapus Semua</button>
            <button className="btn-header-green confirm">✔ Konfirmasi Semua Pesan</button>
        </div>
      </div>

      <div className="log-grid">
        {logs.map((log, index) => (
            <ActivityCard key={index} {...log} />
        ))}
      </div>
    </div>
  );
};

// --- KONTEN KELOLA AKUN ---

// Komponen Card Ringkasan Akun (Reusable)
const AccountSummaryCard = ({ title, value, icon, color }) => (
  <div className="account-summary-card" style={{ backgroundColor: color }}>
    <div className="summary-info">
      <p>{title}</p>
      <h3>{value}</h3>
    </div>
    <div className="summary-icon">{icon}</div>
  </div>
);

const KelolaAkunContent = () => {
  // Data Dummy
  const accountSummaries = [
    { title: 'User', value: '01171', icon: '👤', color: '#fcfcfc' }, // Putih
    { title: 'Cafe', value: '00317', icon: '☕', color: '#fcfcfc' }, // Putih
    { title: 'Event Organizer', value: '00317', icon: '🗓', color: '#fcfcfc' }, // Putih
    { title: 'ADMIN', value: '004', icon: '👨‍💼', color: '#fcfcfc' }, // Putih
  ];

  return (
    <div className="sa-content-padding cafe-theme"> {/* Tetap pakai cafe-theme untuk background coklat */}
      <h1 className="page-title">Kelola Akun</h1>

      {/* Grid Ringkasan Akun */}
      <div className="account-summary-grid">
        {accountSummaries.map((item, index) => (
          <AccountSummaryCard key={index} {...item} />
        ))}
      </div>

      {/* Bagian Greeting & Background Image */}
      <div className="admin-greeting-card">
        <h2>Halo, SUPER ADMIN X Y</h2>
        {/* Pattern Background diatur via CSS */}
      </div>
    </div>
  );
};

// ==========================================
// 3. LAYOUT (SIDEBAR & HEADER)
// ==========================================

const Sidebar = ({ activePage, setActivePage }) => (
  <aside className="sa-sidebar">
    <div className="sa-logo">
      <span className="brand-bold">CARI</span> <span className="brand-light">Spot</span>
    </div>
    <nav className="sa-nav">
      <div 
        className={`nav-item ${activePage === 'dashboard' ? 'active' : ''}`}
        onClick={() => setActivePage('dashboard')}
      >
        <span className="icon">📊</span> Dashboard
      </div>
      <div 
        className={`nav-item ${activePage === 'cafe' ? 'active' : ''}`}
        onClick={() => setActivePage('cafe')}
      >
        <span className="icon">☕</span> Cafe
      </div>
      
      <div 
        className={`nav-item ${activePage === 'event' ? 'active' : ''}`}
        onClick={() => setActivePage('event')}
      >
        <span className="icon">🗓</span> Event
      </div>

      <div 
        className={`nav-item ${activePage === 'log' ? 'active' : ''}`}
        onClick={() => setActivePage('log')}
      >
        <span className="icon">📝</span> Log Aktivitas
      </div>

      <div 
        className={`nav-item ${activePage === 'kelola-akun' ? 'active' : ''}`}
        onClick={() => setActivePage('kelola-akun')}
      >
        <span className="icon">👤</span> Kelola Akun
      </div>
    </nav>

    <div className="sa-nav-bottom">
      <div className="nav-item"><span className="icon">⚙</span> Settings</div>
      <div className="nav-item"><span className="icon">🚪</span> Logout</div>
    </div>
  </aside>
);

const TopHeader = () => (
  <header className="sa-header">
    <div className="search-box">
      <input type="text" placeholder="Cari Cafe/User" />
    </div>
    <div className="header-right">
      <div className="maintenance-badge">🛠 MAINTENANCE SERVER</div>
      <div className="flag-icon">🇮🇩</div>
      <div className="user-profile">
        <div className="avatar-circle">XY</div> 
        <div className="user-text">
            <span className="name">X Y</span>
            <span className="role">Super Admin</span>
        </div>
      </div>
    </div>
  </header>
);

// ==========================================
// 4. KOMPONEN UTAMA (EXPORT)
// ==========================================

const Superadmin = () => {
  const [activePage, setActivePage] = useState('dashboard');

  return (
    <div className="superadmin-layout">
      <Sidebar activePage={activePage} setActivePage={setActivePage} />
      
      <main className="superadmin-main">
        <TopHeader />
        <div className="scrollable-content">
            {activePage === 'dashboard' && <DashboardContent />}
            {activePage === 'cafe' && <CafeContent />}
            {activePage === 'event' && <EventContent />}
            {activePage === 'log' && <LogAktivitasContent />}
            {activePage === 'kelola-akun' && <KelolaAkunContent />}
        </div>
      </main>
    </div>
  );
};

export default Superadmin;