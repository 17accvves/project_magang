import React, { useEffect, useState } from "react";
import "./Menu.css";
import { FaTrash, FaEdit, FaSearch, FaExclamationTriangle, FaSync } from "react-icons/fa";

const API_URL = "http://localhost:8080/menus";
const UPLOAD_URL = "http://localhost:8080/upload";

// Simple fetch dengan logging detail
const apiFetch = async (url, options = {}) => {
  const fullUrl = url.startsWith('http') ? url : `http://localhost:8080${url}`;
  
  console.group('🌐 API Request');
  console.log('URL:', fullUrl);
  console.log('Method:', options.method || 'GET');
  console.log('Headers:', options.headers);
  console.log('Body:', options.body);
  console.groupEnd();

  try {
    const response = await fetch(fullUrl, {
      mode: 'cors',
      credentials: 'omit',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    console.group('📨 API Response');
    console.log('Status:', response.status, response.statusText);
    console.log('Headers:', Object.fromEntries(response.headers.entries()));
    console.groupEnd();

    if (!response.ok) {
      let errorMessage = `HTTP ${response.status}`;
      try {
        const errorData = await response.text();
        errorMessage += ` - ${errorData}`;
      } catch (e) {
        errorMessage += ` - ${response.statusText}`;
      }
      throw new Error(errorMessage);
    }

    // Untuk response tanpa content
    if (response.status === 204) {
      return null;
    }

    const data = await response.json();
    console.log('📊 Response Data:', data);
    return data;

  } catch (error) {
    console.group('❌ API Error');
    console.error('Error:', error);
    console.error('Error name:', error.name);
    console.error('Error message:', error.message);
    console.groupEnd();
    throw error;
  }
};

const Menu = () => {
  const [menus, setMenus] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [showPopup, setShowPopup] = useState(false);
  const [editItem, setEditItem] = useState(null);
  const [formData, setFormData] = useState({
    name: "",
    price: "",
    discount: "",
    startDate: "",
    endDate: "",
    category: "",
    status: "Aktif",
    img: "",
  });
  const [previewImg, setPreviewImg] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    console.log('🔹 Menu component mounted');
    fetchMenus();
  }, []);

  const fetchMenus = async () => {
    try {
      setLoading(true);
      setError("");
      console.log('🔄 Starting fetchMenus...');
      
      const data = await apiFetch(API_URL);
      console.log('✅ fetchMenus success, data:', data);
      setMenus(data || []);
      
    } catch (err) {
      console.error('❌ fetchMenus failed:', err);
      setError(`Gagal mengambil data: ${err.message}`);
      setMenus([]);
    } finally {
      setLoading(false);
    }
  };

  const testConnection = async () => {
    try {
      setError("");
      console.log('🔍 Testing connection to backend...');
      
      const health = await apiFetch('http://localhost:8080/health');
      alert(`✅ Backend Connected!\n\nStatus: ${health.status}\nService: ${health.service}\nTime: ${health.time}`);
      
    } catch (err) {
      const errorMsg = `❌ Cannot connect to backend\n\nPlease check:\n1. Backend server is running\n2. Port 8080 is available\n3. No firewall blocking\n\nError: ${err.message}`;
      setError(errorMsg);
      alert(errorMsg);
    }
  };

  const filteredMenus = menus.filter(
    (m) =>
      m.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      m.category?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const data = new FormData();
    data.append("image", file);

    try {
      console.log('📤 Uploading image...');
      const response = await fetch(UPLOAD_URL, { 
        method: "POST", 
        body: data,
      });
      
      if (!response.ok) {
        throw new Error(`Upload failed: ${response.status}`);
      }
      
      const result = await response.json();
      setFormData({ ...formData, img: result.url });
      setPreviewImg(result.url);
    } catch (err) {
      console.error("Upload gagal:", err);
      alert("Upload gagal: " + err.message);
    }
  };

  const openAddMenu = () => {
    setEditItem(null);
    setFormData({
      name: "",
      price: "",
      discount: "",
      startDate: "",
      endDate: "",
      category: "",
      status: "Aktif",
      img: "",
    });
    setPreviewImg("");
    setShowPopup(true);
  };

  const handleSave = async (e) => {
    e.preventDefault();
    if (!formData.name || !formData.price || !formData.category) {
      alert("Isi semua field wajib!");
      return;
    }

    const price = parseFloat(formData.price);
    const discount = parseFloat(formData.discount || 0);
    const discountedPrice = discount > 0 ? price - (price * discount) / 100 : price;
    
    const newData = { 
      ...formData, 
      price, 
      discount, 
      discountedPrice 
    };

    try {
      console.log('💾 Saving menu...', newData);
      
      if (editItem) {
        await apiFetch(`${API_URL}/${editItem.id}`, {
          method: "PUT",
          body: JSON.stringify(newData),
        });
      } else {
        await apiFetch(API_URL, {
          method: "POST",
          body: JSON.stringify(newData),
        });
      }
      
      setShowPopup(false);
      fetchMenus();
    } catch (err) {
      console.error("Gagal simpan:", err);
      alert("Gagal menyimpan: " + err.message);
    }
  };

  const handleEdit = (menu) => {
    setEditItem(menu);
    setFormData({
      ...menu,
      price: menu.price?.toString() || "",
      discount: menu.discount ? menu.discount.toString() : "",
    });
    setPreviewImg(menu.img || "");
    setShowPopup(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm("Hapus menu ini?")) {
      try {
        await apiFetch(`${API_URL}/${id}`, { method: "DELETE" });
        fetchMenus();
      } catch (err) {
        console.error("Gagal hapus:", err);
        alert("Gagal menghapus: " + err.message);
      }
    }
  };

  const toggleStatus = async (menu) => {
    const newStatus = menu.status === "Aktif" ? "Nonaktif" : "Aktif";
    try {
      await apiFetch(`${API_URL}/${menu.id}`, {
        method: "PUT",
        body: JSON.stringify({ status: newStatus }),
      });
      fetchMenus();
    } catch (err) {
      console.error("Gagal update status:", err);
      alert("Gagal mengubah status: " + err.message);
    }
  };

  return (
    <div className="menu-container">
      <div className="menu-header">
        <h2 className="daftar-title">Daftar Menu</h2>
        <div className="header-actions">
          <button className="test-btn" onClick={testConnection}>
            <FaSync /> Test Koneksi
          </button>
          <button className="refresh-btn" onClick={fetchMenus}>
            Refresh
          </button>
          <button className="add-menu-btn" onClick={openAddMenu}>
            Tambah Menu
          </button>
        </div>
      </div>

      {error && (
        <div className="error-banner">
          <FaExclamationTriangle /> {error}
        </div>
      )}

      <div className="search-container">
        <FaSearch className="search-icon" />
        <input
          type="text"
          placeholder="Cari menu atau kategori..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      {loading && (
        <div className="loading">
          <FaSync className="spinner" /> Memuat data...
        </div>
      )}

      <div className="table-info">
        Menampilkan {filteredMenus.length} dari {menus.length} menu
      </div>

      <table className="menu-table">
        <thead>
          <tr>
            <th>Foto</th>
            <th>Nama</th>
            <th>Harga</th>
            <th>Diskon</th>
            <th>Kategori</th>
            <th>Periode</th>
            <th>Status</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          {filteredMenus.length > 0 ? (
            filteredMenus.map((item) => (
              <tr key={item.id}>
                <td>
                  {item.img ? (
                    <img src={item.img} alt={item.name} className="menu-image" />
                  ) : (
                    <div className="no-image">No Image</div>
                  )}
                </td>
                <td>{item.name}</td>
                <td>
                  {item.discount > 0 ? (
                    <>
                      <span className="old-price">Rp {item.price?.toLocaleString()}</span>
                      <br />
                      <span className="new-price">Rp {item.discountedPrice?.toLocaleString()}</span>
                    </>
                  ) : (
                    `Rp ${item.price?.toLocaleString()}`
                  )}
                </td>
                <td>{item.discount > 0 ? `${item.discount}%` : "-"}</td>
                <td>{item.category}</td>
                <td>
                  {item.discount > 0 
                    ? `${item.startDate || "-"} s/d ${item.endDate || "-"}` 
                    : "-"
                  }
                </td>
                <td>
                  <button
                    className={`status-btn ${item.status === "Aktif" ? "active" : "inactive"}`}
                    onClick={() => toggleStatus(item)}
                  >
                    {item.status}
                  </button>
                </td>
                <td className="aksi-col">
                   <div className="actions">
                     <FaEdit className="edit-icon" onClick={() => handleEdit(item)} title="Edit" />
                     <FaTrash className="delete-icon" onClick={() => handleDelete(item.id)} title="Hapus" />
                   </div>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="8" className="no-result">
                {!loading && menus.length === 0 ? "Belum ada data menu" : "Tidak ada menu yang cocok"}
              </td>
            </tr>
          )}
        </tbody>
      </table>

      {showPopup && (
        <div className="popup-overlay">
          <div className="popup-form">
            <h3>{editItem ? "Edit Menu" : "Tambah Menu"}</h3>
            <form onSubmit={handleSave}>
              <input
                type="text"
                placeholder="Nama Menu *"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
              />
              <input
                type="number"
                placeholder="Harga *"
                value={formData.price}
                onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                min="0"
                step="100"
                required
              />
              <input
                type="number"
                placeholder="Diskon %"
                value={formData.discount}
                onChange={(e) => setFormData({ ...formData, discount: e.target.value })}
                min="0"
                max="100"
                step="1"
              />
              <div className="date-section">
                <label>Periode Diskon:</label>
                <div className="date-inputs">
                  <input
                    type="date"
                    value={formData.startDate}
                    onChange={(e) => setFormData({ ...formData, startDate: e.target.value })}
                  />
                  <input
                    type="date"
                    value={formData.endDate}
                    onChange={(e) => setFormData({ ...formData, endDate: e.target.value })}
                  />
                </div>
              </div>
              <select
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                required
              >
                <option value="">Pilih Kategori *</option>
                <option value="Makanan">Makanan</option>
                <option value="Minuman">Minuman</option>
              </select>
              <div className="upload-section">
                <label>Upload Gambar:</label>
                <input type="file" accept="image/*" onChange={handleImageUpload} />
                {previewImg && <img src={previewImg} alt="Preview" className="preview-image" />}
              </div>
              <div className="popup-buttons">
                <button type="submit" className="save-btn">
                  {editItem ? "Update" : "Simpan"}
                </button>
                <button type="button" className="cancel-btn" onClick={() => setShowPopup(false)}>
                  Batal
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Menu;