import React, { useState } from 'react';
import './SuperAdminLayout.css';

// Import komponen lokal dari folder yang sama
import SASidebar from './SASidebar';
import SATopHeader from './SATopHeader';
import DashboardPage from './DashboardPage';
import CafePage from './CafePage';
import EventPage from './EventPage';
import LogAktivitasPage from './LogAktivitasPage';
import KelolaAkunPage from './KelolaAkunPage';

const SuperAdminLayout = () => {
  const [activePage, setActivePage] = useState('dashboard');

  return (
    <div className="sa-layout">
      <SASidebar activePage={activePage} setActivePage={setActivePage} />
      
      <main className="sa-main">
        <SATopHeader />
        
        <div className="sa-scrollable-content">
            {activePage === 'dashboard' && <DashboardPage />}
            {activePage === 'cafe' && <CafePage />}
            {activePage === 'event' && <EventPage />}
            {activePage === 'log' && <LogAktivitasPage />}
            {activePage === 'kelola-akun' && <KelolaAkunPage />}
        </div>
      </main>
    </div>
  );
};

export default SuperAdminLayout;