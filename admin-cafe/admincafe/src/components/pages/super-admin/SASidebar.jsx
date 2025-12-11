import React from 'react';
import './SASidebar.css';

const SASidebar = ({ activePage, setActivePage }) => {
  return (
    <aside className="sa-sidebar-container">
      <div className="sa-logo">
        <span className="brand-bold">CARI</span> <span className="brand-light">Spot</span>
      </div>

      <nav className="sa-nav-list">
        <div 
            className={`sa-nav-item ${activePage === 'dashboard' ? 'active' : ''}`}
            onClick={() => setActivePage('dashboard')}
        >
            <span className="icon">📊</span> Dashboard
        </div>
        <div 
            className={`sa-nav-item ${activePage === 'cafe' ? 'active' : ''}`}
            onClick={() => setActivePage('cafe')}
        >
            <span className="icon">☕</span> Cafe
        </div>
        <div 
            className={`sa-nav-item ${activePage === 'event' ? 'active' : ''}`}
            onClick={() => setActivePage('event')}
        >
            <span className="icon">🗓️</span> Event
        </div>
        <div 
            className={`sa-nav-item ${activePage === 'log' ? 'active' : ''}`}
            onClick={() => setActivePage('log')}
        >
            <span className="icon">📝</span> Log Aktivitas
        </div>
        <div 
            className={`sa-nav-item ${activePage === 'kelola-akun' ? 'active' : ''}`}
            onClick={() => setActivePage('kelola-akun')}
        >
            <span className="icon">👤</span> Kelola Akun
        </div>
      </nav>

      <div className="sa-nav-bottom">
        <div className="sa-nav-item"><span className="icon">⚙️</span> Settings</div>
        <div className="sa-nav-item"><span className="icon">🚪</span> Logout</div>
      </div>
    </aside>
  );
};

export default SASidebar;