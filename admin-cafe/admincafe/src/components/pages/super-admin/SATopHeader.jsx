import React from 'react';
import './SATopHeader.css';

const SATopHeader = () => {
  return (
    <header className="sa-top-header">
      <div className="sa-search-box">
        <input type="text" placeholder="Cari Cafe/User" />
      </div>
      <div className="sa-header-right">
        <div className="sa-maintenance-badge">🛠️ MAINTENANCE SERVER</div>
        <div className="sa-flag-icon">🇮🇩</div>
        <div className="sa-user-profile">
          <div className="sa-avatar-circle">XY</div> 
          <div className="sa-user-text">
              <span className="name">X Y</span>
              <span className="role">Super Admin</span>
          </div>
        </div>
      </div>
    </header>
  );
};

export default SATopHeader;