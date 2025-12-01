import React from 'react';
import './DashboardPage.css';

const StatCard = ({ title, value, icon, colorClass }) => (
    <div className={`d-stat-card ${colorClass}`}>
        <div className="stat-info">
            <p>{title}</p>
            <h2>{value}</h2>
        </div>
        <div className="stat-icon">{icon}</div>
    </div>
);

const DashboardPage = () => {
    const stats = [
        { title: 'User', value: '01171', icon: '👤', colorClass: 'bg-purple' },
        { title: 'Kafe Terdaftar', value: '00317', icon: '☕', colorClass: 'bg-orange' },
        { title: 'Event Organizer', value: '00317', icon: '📅', colorClass: 'bg-blue' },
        { title: 'Admin Aktif', value: '004', icon: '👮', colorClass: 'bg-green' },
    ];

    return (
        <div className="d-content-padding">
            <h1 className="d-page-title">Dashboard</h1>
            
            <div className="d-grid-stats">
                {stats.map((s, i) => <StatCard key={i} {...s} />)}
            </div>

            <div className="d-grid-main">
                <div className="d-card chart-section">
                    <div className="d-card-header">
                        <h3>Pengunjung</h3>
                        <button className="d-btn-blue">Lihat</button>
                    </div>
                    <div className="d-chart-placeholder">
                        <div className="d-chart-lines">📈 [Grafik Placeholder]</div>
                        <div className="d-chart-legend">
                            <span>🟣 Loyal</span> <span>🔴 Baru</span> <span>🟢 Premium</span>
                        </div>
                    </div>
                </div>

                <div className="d-card calendar-section">
                    <h3>Sep, 19 Friday</h3>
                    <div className="d-calendar-placeholder">
                       <div className="d-cal-grid">
                           <span>M</span><span>T</span><span>W</span><span>T</span><span>F</span><span className="red">S</span><span className="red">S</span>
                           <span>18</span><span className="today">19</span><span>20</span><span>21</span><span>22</span><span className="red">23</span><span className="red">24</span>
                       </div>
                    </div>
                </div>
            </div>

            <div className="d-card table-section">
                <div className="d-card-header">
                    <h3>CAFE TERFAVORIT</h3>
                    <select><option>September</option></select>
                </div>
                <table className="d-table">
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
                            <td><span className="d-badge active">Aktif</span></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default DashboardPage;