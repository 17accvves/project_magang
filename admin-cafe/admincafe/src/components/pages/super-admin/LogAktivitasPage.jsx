import React from 'react';
import './LogAktivitasPage.css';

const ActivityCard = ({ title, date, status, subject, claim, variant }) => (
    <div className={`activity-card ${variant}`}>
      <div className="act-content">
        <h3 className="act-title">{title}</h3>
        <div className="act-meta">
            <span className="act-date">{date}</span>
            <div className="act-status">Operasional : <span className="status-pill">{status}</span></div>
        </div>
        <div className="act-details">
            <p>Perihal : {subject}</p>
            <p>{claim}</p>
        </div>
      </div>
      <div className="act-actions">
        <button className="btn-act confirm"><span className="icon-check">✔</span> Konfirmasi</button>
        <button className="btn-act history"><span className="icon-dots">•••</span> Riwayat</button>
        <button className="btn-act chat"><span className="icon-chat">💬</span> Buka Chat</button>
      </div>
    </div>
);

const LogAktivitasPage = () => {
  const logs = [
    { title: 'Kafe xxx mengirimkan perubahan Diskon', date: '7 Sep 2025', status: 'Aktif', subject: 'Diskon 20% sd 20rb', claim: 'Klaim : Anggota Premium', variant: 'dark' },
    { title: 'Kafe xxx mengirimkan perubahan Lokasi', date: '7 Sep 2025', status: 'Aktif', subject: 'Perubahan Lokasi', claim: 'Klaim : SUPER ADMIN', variant: 'light' },
    { title: 'Kafe xxx melakukan perubahan foto profil', date: '7 Sep 2025', status: 'Aktif', subject: 'Perubahan Akun', claim: '', variant: 'light' },
    { title: 'USER Delon berlangganan PREMIUM', date: '7 Sep 2025', status: 'Aktif', subject: 'Berlangganan', claim: 'Klaim : Anggota Premium', variant: 'dark' },
    { title: 'Kafe xxx mengirimkan perubahan Diskon', date: '7 Sep 2025', status: 'Aktif', subject: 'Diskon 20% sd 20rb', claim: 'Klaim : Anggota Premium', variant: 'dark' },
    { title: 'Kafe xxx mengirimkan perubahan Lokasi', date: '7 Sep 2025', status: 'Aktif', subject: 'Perubahan Lokasi', claim: 'Datang : 75', variant: 'light' },
  ];

  return (
    <div className="cafe-content-padding cafe-theme">
      <div className="log-header">
        <h1 className="page-title">Log Aktivitas</h1>
        <div className="log-header-actions">
            <button className="btn-header-green delete">🗑 Hapus Semua</button>
            <button className="btn-header-green confirm">✔ Konfirmasi Semua Pesan</button>
        </div>
      </div>
      <div className="log-grid">
        {logs.map((log, index) => <ActivityCard key={index} {...log} />)}
      </div>
    </div>
  );
};

export default LogAktivitasPage;