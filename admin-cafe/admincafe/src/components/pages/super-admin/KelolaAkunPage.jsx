import React from 'react';
import './KelolaAkunPage.css';

const AccountSummaryCard = ({ title, value, icon, color }) => (
  <div className="account-summary-card" style={{ backgroundColor: color }}>
    <div className="summary-info">
      <p>{title}</p>
      <h3>{value}</h3>
    </div>
    <div className="summary-icon">{icon}</div>
  </div>
);

const KelolaAkunPage = () => {
  const accountSummaries = [
    { title: 'User', value: '01171', icon: '👤', color: '#fcfcfc' },
    { title: 'Cafe', value: '00317', icon: '☕', color: '#fcfcfc' },
    { title: 'Event Organizer', value: '00317', icon: '🗓️', color: '#fcfcfc' },
    { title: 'ADMIN', value: '004', icon: '👨‍💼', color: '#fcfcfc' },
  ];

  return (
    <div className="cafe-content-padding cafe-theme">
      <h1 className="page-title">Kelola Akun</h1>
      <div className="account-summary-grid">
        {accountSummaries.map((item, index) => <AccountSummaryCard key={index} {...item} />)}
      </div>
      <div className="admin-greeting-card">
        <h2>Halo, SUPER ADMIN X Y</h2>
      </div>
    </div>
  );
};

export default KelolaAkunPage;