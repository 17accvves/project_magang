import React from 'react';
import './EventPage.css';

const EventPage = () => {
  const runningEvents = [
    { id: 1, cafe: 'Yotta.id', desc: 'band : Avenged 7x', eventName: 'Kelayapan', organizer: 'SMA Teykom', date: '19-09-2025', time: '17.00 WITA-Selesai', type: 'music' },
    { id: 2, cafe: 'Ngopisantuy', desc: 'band : UNM Band', eventName: 'Jam Malam', organizer: 'SMA Teykom', date: '19-09-2025', time: '14.00 WITA-Selesai', type: 'coffee' },
    { id: 3, cafe: 'Teh Indonesia', desc: 'band : free to all', eventName: 'spot.in', organizer: 'SMA Teykom', date: '19-09-2025', time: '22.00 WITA-Selesai', type: 'drink' },
    { id: 4, cafe: 'Yotta.id', desc: 'band : Avenged 7x', eventName: 'Kelayapan', organizer: 'SMA Teykom', date: '19-09-2025', time: '17.00 WITA-Selesai', type: 'music' },
    { id: 5, cafe: 'Ngopisantuy', desc: 'band : UNM Band', eventName: 'Jam Malam', organizer: 'SMA Teykom', date: '19-09-2025', time: '14.00 WITA-Selesai', type: 'coffee' },
    { id: 6, cafe: 'Teh Indonesia', desc: 'band : free to all', eventName: 'spot.in', organizer: 'SMA Teykom', date: '19-09-2025', time: '22.00 WITA-Selesai', type: 'drink' },
  ];

  return (
    <div className="cafe-content-padding cafe-theme">
      <h1 className="page-title">Event</h1>

      <div className="event-layout-grid">
        {/* Kolom Kiri */}
        <div className="event-left-col">
          <div className="sa-card eo-active-card">
            <div className="eo-info">
              <p className="label">EO Aktif</p>
              <h2 className="value">00074</h2>
            </div>
            <button className="btn-view-eo">📄 Lihat EO</button>
          </div>

          <div className="sa-card running-event-card">
            <div className="card-header-center">
                <h3>Event Berjalan</h3>
                <div className="scroll-arrow up">︿</div>
            </div>

            <div className="event-table-header">
                <div className="col-header"><span className="header-icon">🎵</span> Live Music</div>
                <div className="col-header"><span className="header-icon">🍔</span> Bazar</div>
                <div className="col-header"><span className="header-icon">🕒</span> Waktu</div>
            </div>

            <div className="event-list-container">
                {runningEvents.map((item, index) => (
                    <div className="event-row" key={index}>
                        <div className="event-col col-1">
                            <span className="row-num">{item.id}</span>
                            <div className="col-icon type-icon">{item.type === 'music' ? '🥤' : item.type === 'coffee' ? '☕' : '🥤'}</div>
                            <div className="col-text"><strong>{item.cafe}</strong><span>{item.desc}</span></div>
                        </div>
                        <div className="event-col col-2">
                            <div className="col-icon">🛵</div>
                            <div className="col-text"><strong>{item.eventName}</strong><span>penyelenggara : {item.organizer}</span></div>
                        </div>
                        <div className="event-col col-3">
                            <div className="col-text align-right"><strong>{item.date}</strong><span>Pukul : {item.time}</span></div>
                        </div>
                    </div>
                ))}
            </div>
            <div className="card-footer-center"><div className="scroll-arrow down">﹀</div></div>
          </div>
        </div>

        {/* Kolom Kanan */}
        <div className="event-right-col">
          <div className="sa-card registration-card">
            <h3>Registrasi</h3>
            <div className="reg-buttons">
                <button className="btn-reg blue">Verifikasi Akun EO</button>
                <button className="btn-reg blue">Kelola Akun EO</button>
            </div>
          </div>

          <div className="sa-card manage-event-card">
            <h3>Manage Event</h3>
            <div className="manage-item"><span>Edit</span><button className="btn-icon">✏️</button></div>
            <div className="manage-item"><span>Lihat Riwayat</span><button className="btn-icon">↺</button></div>
          </div>

          <div className="sa-card highlight-card">
            <div className="highlight-img">
                <img src="https://via.placeholder.com/300x150" alt="Event" />
                <div className="tag-yellow">Yotta.id</div>
            </div>
            <div className="highlight-info">
                <div className="date-box"><span className="month">AUG</span><span className="date">27</span></div>
                <div className="detail-box">
                    <h4>Bazar SMAN 1 Kilo</h4>
                    <p>Outlet Pettarani, Makassar</p>
                    <p className="time">14:20 WITA</p>
                    <div className="attendees">👥 34 • ⭐⭐⭐⭐⭐</div>
                </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EventPage;