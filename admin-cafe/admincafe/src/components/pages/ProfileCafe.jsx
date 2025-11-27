// ProfileCafe.jsx
import React, { useState, useEffect } from "react";
import {
  FaTrash,
  FaInstagram,
  FaFacebook,
  FaTiktok,
  FaPhone,
  FaClock,
  FaTimes,
  FaPlus,
  FaCamera,
  FaExclamationTriangle,
  FaBug
} from "react-icons/fa";
import "./ProfileCafe.css";

// Fungsi helper untuk mendapatkan full image URL
const getFullImageUrl = (url) => {
  if (!url) return "../../src/assets/maincafe.jpg";
  
  if (url.startsWith('http')) return url;
  return `http://localhost:8080${url}`;
};

// Fungsi untuk convert ISO time ke format HH:MM
const convertISOToTime = (isoTime) => {
  if (!isoTime) return "08:00";
  
  // Jika sudah format HH:MM, return langsung
  if (isoTime.match(/^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/)) {
    return isoTime;
  }
  
  // Jika format ISO (0000-01-01T19:18:00Z)
  if (isoTime.includes('T')) {
    try {
      const date = new Date(isoTime);
      const hours = date.getHours().toString().padStart(2, '0');
      const minutes = date.getMinutes().toString().padStart(2, '0');
      return `${hours}:${minutes}`;
    } catch (error) {
      console.error('Error converting ISO time:', error);
      return "08:00";
    }
  }
  
  return "08:00";
};

// Fungsi untuk test load gambar
const testImageLoad = (url) => {
  return new Promise((resolve) => {
    const testImg = new Image();
    testImg.onload = () => resolve(true);
    testImg.onerror = () => resolve(false);
    testImg.src = url;
  });
};

function ProfileCafe() {
  const [gallery, setGallery] = useState([]);
  const [showEdit, setShowEdit] = useState(false);
  const [statusBuka, setStatusBuka] = useState({
    status: "",
    waktu: "",
    warna: "#ff7878",
    hari: ""
  });
  const [isSaving, setIsSaving] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [maxGalleryReached, setMaxGalleryReached] = useState(false);

  const [formData, setFormData] = useState({
    id: 1,
    nama: "",
    alamat: "",
    telepon: "",
    deskripsi: "",
    main_image: "",
    verified: false,
    social_media: [],
    operational_hours: [],
    facilities: []
  });

  // === FUNGSI DEBUG UNTUK MEMECAHKAN MASALAH JAM ===
  const debugOperationalHours = async () => {
    try {
      console.log('🔍 === DEBUG OPERATIONAL HOURS ===');
      
      const dbResponse = await fetch('http://localhost:8080/cafe/operational-hours');
      if (dbResponse.ok) {
        const dbData = await dbResponse.json();
        console.log('🗄️ Database operational hours (RAW):', dbData);
        
        console.log('📅 Detailed database hours:');
        dbData.forEach(hour => {
          console.log(`   ${hour.hari}: "${hour.buka}" - "${hour.tutup}"`);
          console.log(`     Converted: ${convertISOToTime(hour.buka)} - ${convertISOToTime(hour.tutup)}`);
        });
      }
      
      console.log('📝 FormData operational_hours:', formData.operational_hours);
      console.log('🎯 statusBuka:', statusBuka);

      // Debug status saat ini
      const now = new Date();
      const currentTime = now.getHours() * 60 + now.getMinutes();
      console.log('🕒 Current time:', {
        hours: now.getHours(),
        minutes: now.getMinutes(),
        totalMinutes: currentTime,
        formatted: `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
      });
      
      console.log('🔍 === END DEBUG ===');
      
    } catch (error) {
      console.error('❌ Debug error:', error);
    }
  };

  // === FETCH DATA DARI BACKEND ===
  useEffect(() => {
    fetchCafeProfile();
    fetchGallery();
  }, []);

  useEffect(() => {
    setMaxGalleryReached(gallery.length >= 10);
  }, [gallery]);

  // === LOGIKA CEK HARI DAN JAM BUKA YANG DIPERBAIKI ===
  useEffect(() => {
    const checkStatus = () => {
      const hariList = ["Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"];
      const hariListLower = ["minggu", "senin", "selasa", "rabu", "kamis", "jumat", "sabtu"];
      const now = new Date();
      const dayIndex = now.getDay();
      const currentDay = hariListLower[dayIndex];
      const currentTime = now.getHours() * 60 + now.getMinutes();

      console.log(`🕒 Checking status for ${currentDay} at ${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`);

      const todaySchedule = formData.operational_hours.find(
        hour => hour.hari.toLowerCase() === currentDay
      );

      if (!todaySchedule) {
        setStatusBuka({
          status: "Tidak ada jadwal",
          waktu: "Belum diatur",
          warna: "#ff7878",
          hari: hariList[dayIndex]
        });
        return;
      }

      const buka = convertISOToTime(todaySchedule.buka);
      const tutup = convertISOToTime(todaySchedule.tutup);
      
      const [bukaH, bukaM] = buka.split(":").map(Number);
      const [tutupH, tutupM] = tutup.split(":").map(Number);
      
      const bukaMinutes = bukaH * 60 + bukaM;
      const tutupMinutes = tutupH * 60 + tutupM;

      // Debug detail
      console.log('📋 Schedule details:', {
        hari: currentDay,
        buka,
        tutup,
        bukaMinutes,
        tutupMinutes,
        currentTime,
        currentTimeFormatted: `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
      });

      let isOpen = false;

      // Logic untuk handle buka sampai lewat tengah malam
      if (tutupMinutes < bukaMinutes) {
        // Contoh: 19:00 - 02:00 (buka malam sampai pagi)
        isOpen = currentTime >= bukaMinutes || currentTime <= tutupMinutes;
        console.log('🌙 Buka sampai lewat tengah malam');
      } else {
        // Contoh: 08:00 - 17:00 (buka siang)
        isOpen = currentTime >= bukaMinutes && currentTime <= tutupMinutes;
        console.log('☀️ Buka dalam hari yang sama');
      }

      console.log(`🔓 Status: ${isOpen ? 'BUKA' : 'TUTUP'}`);

      setStatusBuka({
        status: isOpen ? "Buka Sekarang" : "Tutup Sekarang",
        waktu: `${buka} - ${tutup}`,
        warna: isOpen ? "#00c896" : "#ff7878",
        hari: hariList[dayIndex]
      });
    };

    checkStatus();
    // Update setiap 30 detik (lebih responsif)
    const interval = setInterval(checkStatus, 30000);
    return () => clearInterval(interval);
  }, [formData.operational_hours]);

  const fetchCafeProfile = async () => {
    try {
      console.log('🔄 Fetching cafe profile...');
      setIsLoading(true);
      
      const response = await fetch('http://localhost:8080/cafe/profile');
      console.log('📡 Response status:', response.status);
      
      if (response.ok) {
        const data = await response.json();
        console.log('✅ Cafe profile fetched:', data);
        
        // Jika data kosong, gunakan default values
        if (!data.nama) {
          console.log('ℹ️ No profile data, using defaults');
          setFormData({
            id: 1,
            nama: "",
            alamat: "",
            telepon: "",
            deskripsi: "",
            main_image: "",
            verified: false,
            social_media: [
              { platform: "instagram", url: "" },
              { platform: "facebook", url: "" },
              { platform: "tiktok", url: "" }
            ],
            operational_hours: [
              { hari: "senin", buka: "08:00", tutup: "17:00" },
              { hari: "selasa", buka: "08:00", tutup: "17:00" },
              { hari: "rabu", buka: "08:00", tutup: "17:00" },
              { hari: "kamis", buka: "08:00", tutup: "17:00" },
              { hari: "jumat", buka: "08:00", tutup: "17:00" },
              { hari: "sabtu", buka: "08:00", tutup: "17:00" },
              { hari: "minggu", buka: "08:00", tutup: "17:00" }
            ],
            facilities: [
            { nama_fasilitas: "wifi", tersedia: false },
            { nama_fasilitas: "smoking_area", tersedia: false},
            { nama_fasilitas: "indoor", tersedia: false},
            { nama_fasilitas: "parkiran", tersedia: false},
            { nama_fasilitas: "musholla", tersedia: false},
            { nama_fasilitas: "spot_foto", tersedia: false},
            { nama_fasilitas: "live_music", tersedia: false },
            { nama_fasilitas: "outdoor", tersedia: false },
            { nama_fasilitas: "air_conditioner", tersedia: false }
            ]
          });
        } else {
          console.log('ℹ️ Profile data found, setting form data');
          
          // PROCESS OPERATIONAL HOURS - CONVERT ISO TO SIMPLE TIME
          let operationalHours = [];
          if (data.operational_hours && data.operational_hours.length > 0) {
            operationalHours = data.operational_hours.map(hour => ({
              hari: hour.hari,
              buka: convertISOToTime(hour.buka),
              tutup: convertISOToTime(hour.tutup)
            }));
            console.log('Converted operational hours:', operationalHours);
          } else {
            operationalHours = [
              { hari: "senin", buka: "08:00", tutup: "17:00" },
              { hari: "selasa", buka: "08:00", tutup: "17:00" },
              { hari: "rabu", buka: "08:00", tutup: "17:00" },
              { hari: "kamis", buka: "08:00", tutup: "17:00" },
              { hari: "jumat", buka: "08:00", tutup: "17:00" },
              { hari: "sabtu", buka: "08:00", tutup: "17:00" },
              { hari: "minggu", buka: "08:00", tutup: "17:00" }
            ];
          }

          setFormData(prev => ({
            ...prev,
            ...data,
            social_media: data.social_media || [
              { platform: "instagram", url: "" },
              { platform: "facebook", url: "" },
              { platform: "tiktok", url: "" }
            ],
            operational_hours: operationalHours,
            facilities: data.facilities || [
            { nama_fasilitas: "wifi", tersedia: false },
            { nama_fasilitas: "smoking_area", tersedia: false},
            { nama_fasilitas: "indoor", tersedia: false},
            { nama_fasilitas: "parkiran", tersedia: false},
            { nama_fasilitas: "musholla", tersedia: false},
            { nama_fasilitas: "spot_foto", tersedia: false},
            { nama_fasilitas: "live_music", tersedia: false },
            { nama_fasilitas: "outdoor", tersedia: false },
            { nama_fasilitas: "air_conditioner", tersedia: false }
            ]
          }));

          // Debug setelah set form data
          console.log('🔍 Debug operational hours setelah set:', {
            operationalHours: operationalHours,
            hariIni: new Date().getDay(),
            hariIniNama: ["minggu", "senin", "selasa", "rabu", "kamis", "jumat", "sabtu"][new Date().getDay()],
            scheduleHariIni: operationalHours.find(hour => 
              hour.hari === ["minggu", "senin", "selasa", "rabu", "kamis", "jumat", "sabtu"][new Date().getDay()]
            )
          });
        }
      } else {
        console.error('❌ Failed to fetch cafe profile');
        setFormData(prev => ({
          ...prev,
          social_media: [
            { platform: "instagram", url: "" },
            { platform: "facebook", url: "" },
            { platform: "tiktok", url: "" }
          ],
          operational_hours: [
            { hari: "senin", buka: "08:00", tutup: "17:00" },
            { hari: "selasa", buka: "08:00", tutup: "17:00" },
            { hari: "rabu", buka: "08:00", tutup: "17:00" },
            { hari: "kamis", buka: "08:00", tutup: "17:00" },
            { hari: "jumat", buka: "08:00", tutup: "17:00" },
            { hari: "sabtu", buka: "08:00", tutup: "17:00" },
            { hari: "minggu", buka: "08:00", tutup: "17:00" }
          ],
          facilities: [
            { nama_fasilitas: "wifi", tersedia: false },
            { nama_fasilitas: "smoking_area", tersedia: false},
            { nama_fasilitas: "indoor", tersedia: false},
            { nama_fasilitas: "parkiran", tersedia: false},
            { nama_fasilitas: "musholla", tersedia: false},
            { nama_fasilitas: "spot_foto", tersedia: false},
            { nama_fasilitas: "live_music", tersedia: false },
            { nama_fasilitas: "outdoor", tersedia: false },
            { nama_fasilitas: "air_conditioner", tersedia: false }
          ]
        }));
      }
    } catch (error) {
      console.error('💥 Error fetching cafe profile:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchGallery = async () => {
    try {
      const response = await fetch('http://localhost:8080/cafe/gallery');
      if (response.ok) {
        const data = await response.json();
        setGallery(data.map(item => ({
          id: item.id,
          src: item.image_url,
          urutan: item.urutan
        })));
      }
    } catch (error) {
      console.error('Error fetching gallery:', error);
    }
  };

  // === FUNGSI UNTUK GAMBAR PROFILE ===
  const handleChangeProfileImage = async (e) => {
    const file = e.target.files[0];
    if (file) {
      try {
        const reader = new FileReader();
        reader.onload = (event) => {
          setFormData(prev => ({
            ...prev,
            main_image: event.target.result
          }));
        };
        reader.readAsDataURL(file);

        const formData = new FormData();
        formData.append('image', file);
        
        const response = await fetch('http://localhost:8080/cafe/profile/image', {
          method: 'POST',
          body: formData
        });
        
        if (response.ok) {
          const result = await response.json();
          setFormData(prev => ({
            ...prev,
            main_image: result.image_url
          }));
          await testImageLoad(result.image_url);
          alert('Gambar profile berhasil diupload!');
        } else {
          const errorText = await response.text();
          alert('Gagal upload gambar profile: ' + errorText);
        }
      } catch (error) {
        alert('Error uploading gambar profile: ' + error.message);
      }
    }
  };

  const handleDeleteProfileImage = async () => {
    if (window.confirm('Apakah Anda yakin ingin menghapus gambar profile?')) {
      try {
        const response = await fetch('http://localhost:8080/cafe/profile', {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            ...formData,
            main_image: ""
          })
        });
        
        if (response.ok) {
          setFormData(prev => ({
            ...prev,
            main_image: ""
          }));
          alert('Gambar profile berhasil dihapus!');
        }
      } catch (error) {
        console.error('Error deleting profile image:', error);
      }
    }
  };

  // === FUNGSI HANDLER LAINNYA ===
  const handleDeleteImage = async (id) => {
    try {
      const response = await fetch(`http://localhost:8080/cafe/gallery/${id}`, {
        method: 'DELETE'
      });
      if (response.ok) {
        setGallery(gallery.filter((img) => img.id !== id));
        alert('Gambar berhasil dihapus dari galeri!');
      } else {
        alert('Gagal menghapus gambar!');
      }
    } catch (error) {
      alert('Error menghapus gambar: ' + error.message);
    }
  };

  const handleAddImage = async (e) => {
    const file = e.target.files[0];
    e.target.value = '';
    
    if (!file) return;

    if (gallery.length >= 10) {
      alert('❌ Maaf, galeri sudah penuh! Maksimal 10 gambar.');
      return;
    }

    try {
      const uploadFormData = new FormData();
      uploadFormData.append('image', file);
      uploadFormData.append('urutan', gallery.length + 1);

      const response = await fetch('http://localhost:8080/cafe/gallery', {
        method: 'POST',
        body: uploadFormData
      });
      
      if (response.ok) {
        const newImage = await response.json();
        setGallery([...gallery, {
          id: newImage.id,
          src: newImage.image_url,
          urutan: newImage.urutan
        }]);
        alert('Gambar berhasil ditambahkan ke galeri!');
      } else {
        alert('Gagal menambahkan gambar ke galeri!');
      }
    } catch (error) {
      alert('Error menambahkan gambar: ' + error.message);
    }
  };

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;

    if (name.startsWith('facility_')) {
      const facilityName = name.replace('facility_', '');
      setFormData({
        ...formData,
        facilities: formData.facilities.map(fac =>
          fac.nama_fasilitas === facilityName 
            ? { ...fac, tersedia: checked }
            : fac
        )
      });
    } else if (name.startsWith('social_')) {
      const platform = name.replace('social_', '');
      setFormData({
        ...formData,
        social_media: formData.social_media.map(sm =>
          sm.platform === platform
            ? { ...sm, url: value }
            : sm
        )
      });
    } else if (name.startsWith('hour_')) {
      const [day, field] = name.replace('hour_', '').split('_');
      
      setFormData({
        ...formData,
        operational_hours: formData.operational_hours.map(hour =>
          hour.hari === day
            ? { ...hour, [field]: value }
            : hour
        )
      });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  // === FUNGSI VALIDASI DATA ===
  const validateDataBeforeSave = () => {
    const errors = [];

    if (!formData.nama || formData.nama.trim() === '') {
      errors.push('Nama cafe tidak boleh kosong');
    }

    formData.operational_hours.forEach(hour => {
      if (!hour.buka || hour.buka.trim() === '') {
        errors.push(`Jam buka untuk hari ${hour.hari} tidak boleh kosong`);
      }
      if (!hour.tutup || hour.tutup.trim() === '') {
        errors.push(`Jam tutup untuk hari ${hour.hari} tidak boleh kosong`);
      }
      
      const timeRegex = /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/;
      if (hour.buka && !timeRegex.test(hour.buka)) {
        errors.push(`Format jam buka ${hour.hari} tidak valid`);
      }
      if (hour.tutup && !timeRegex.test(hour.tutup)) {
        errors.push(`Format jam tutup ${hour.hari} tidak valid`);
      }
    });

    return errors;
  };

  // === FUNGSI UPDATE SOCIAL MEDIA ===
  const updateSocialMedia = async () => {
    try {
      await fetch('http://localhost:8080/cafe/social-media/all', {
        method: 'DELETE'
      });

      for (const sm of formData.social_media) {
        await fetch('http://localhost:8080/cafe/social-media', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            platform: sm.platform,
            url: sm.url
          })
        });
      }
    } catch (error) {
      console.error('Error updating social media:', error);
    }
  };

  // === FUNGSI UPDATE OPERATIONAL HOURS YANG DIPERBAIKI ===
  const updateOperationalHours = async () => {
    try {
      console.log('🔄 Starting operational hours update...');
      
      const operationalData = formData.operational_hours.map(hour => {
        const buka = hour.buka ? hour.buka.padStart(5, '0') : "08:00";
        const tutup = hour.tutup ? hour.tutup.padStart(5, '0') : "17:00";
        
        return {
          hari: hour.hari,
          buka: buka,
          tutup: tutup
        };
      });

      console.log('📤 Final data to send:', operationalData);

      const response = await fetch('http://localhost:8080/cafe/operational-hours', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(operationalData)
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP ${response.status}: ${errorText}`);
      }
      
      console.log('✅ Operational hours updated successfully');
      
    } catch (error) {
      console.error('💥 Error updating operational hours:', error);
      throw error;
    }
  };

  // === FUNGSI UPDATE FACILITIES ===
  const updateFacilities = async () => {
    try {
      const response = await fetch('http://localhost:8080/cafe/facilities', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData.facilities)
      });
      
      if (!response.ok) {
        throw new Error('Gagal menyimpan fasilitas');
      }
    } catch (error) {
      console.error('Error updating facilities:', error);
      throw error;
    }
  };

  // === FUNGSI SIMPAN YANG LEBIH ROBUST ===
  const handleSave = async () => {
    const validationErrors = validateDataBeforeSave();
    if (validationErrors.length > 0) {
      alert('❌ Terdapat kesalahan dalam data:\n' + validationErrors.join('\n'));
      return;
    }

    setIsSaving(true);
    console.log('💾 Starting save process...');
    
    try {
      // Simpan profile utama
      const profileResponse = await fetch('http://localhost:8080/cafe/profile', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: formData.id,
          nama: formData.nama,
          alamat: formData.alamat,
          telepon: formData.telepon,
          deskripsi: formData.deskripsi,
          main_image: formData.main_image,
          verified: formData.verified
        })
      });

      if (!profileResponse.ok) {
        const errorText = await profileResponse.text();
        throw new Error(`Gagal menyimpan profile: ${errorText}`);
      }

      // Update social media
      await updateSocialMedia();

      // Update operational hours
      await updateOperationalHours();

      // Update facilities
      await updateFacilities();

      console.log("🎉 Semua data berhasil disimpan!");
      setShowEdit(false);
      alert('Profile berhasil disimpan!');
      
      // Refresh data
      await fetchCafeProfile();
      
    } catch (error) {
      console.error('❌ Error saving profile:', error);
      alert('Gagal menyimpan profile: ' + error.message);
    } finally {
      setIsSaving(false);
    }
  };

  // === Buka link sosial media ===
  const openSocialLink = (platform) => {
    const social = formData.social_media.find(sm => sm.platform === platform);
    if (!social || !social.url || social.url.trim() === "") {
      alert(`Link ${platform} belum diatur!`);
      return;
    }

    const finalLink = social.url.startsWith("http")
      ? social.url
      : `https://${social.url}`;
    window.open(finalLink, "_blank");
  };

  // === GET FACILITIES TEXT ===
  const getFacilitiesText = () => {
    const availableFacilities = formData.facilities.filter(fac => fac.tersedia);
    
    if (availableFacilities.length === 0) {
      return "Belum ada fasilitas yang dipilih";
    }
    
    return availableFacilities.map(fac => {
      const names = {
        wifi: "Wifi",
        smoking_area: "Smoking Area",
        indoor: "Indoor",
        parkiran: "Parkiran",
        musholla: "Musholla",
        spot_foto: "Spot Foto",
        live_music: "Live Music", 
        outdoor: "Outdoor",
        air_conditioner: "AC"
      };
      return names[fac.nama_fasilitas] || fac.nama_fasilitas;
    }).join(" ✔ ");
  };

  // === RENDER LOADING ===
  if (isLoading) {
    return (
      <div className="profile-cafe-container">
        <div className="loading-container">
          <h2>Memuat Profile Cafe...</h2>
          <div className="loading-spinner"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="profile-cafe-container">
      <h2 className="profile-title">Profile Cafe</h2>

      {/* Bagian Atas */}
      <div className="top-section">
        <div className="main-image-card">
          <div className="image-container">
            <label htmlFor="profileImageInput" className="image-upload-label">
              <img
                src={getFullImageUrl(formData.main_image)}
                alt="Main Cafe"
                className="main-image"
                onError={(e) => {
                  e.target.src = "../../src/assets/LOGO.jpg";
                }}
              />
              <div className="image-overlay">
                <FaCamera className="camera-icon" />
                <span>Klik untuk ganti gambar</span>
              </div>
            </label>
            <input
              id="profileImageInput"
              type="file"
              accept="image/*"
              style={{ display: "none" }}
              onChange={handleChangeProfileImage}
            />
            
            {formData.main_image && formData.main_image !== "" && (
              <button 
                className="delete-profile-btn"
                onClick={handleDeleteProfileImage}
                title="Hapus gambar profile"
              >
                <FaTrash />
              </button>
            )}
          </div>
        </div>

        <div className="info-card">
          <div className="info-row">
            <span className="label">Nama :</span> 
            <span>{formData.nama || "Belum diisi"}</span>
          </div>
          <div className="info-row">
            <span className="label">Alamat :</span>{" "}
            <span>{formData.alamat || "Belum diisi"}</span>
          </div>

          <div className="info-row">
            <span className="label">Sosial Media :</span>
            <span className="social-icons">
              <FaInstagram
                className="social-icon"
                onClick={() => openSocialLink("instagram")}
              />
              <FaTiktok
                className="social-icon"
                onClick={() => openSocialLink("tiktok")}
              />
              <FaFacebook
                className="social-icon"
                onClick={() => openSocialLink("facebook")}
              />
            </span>
          </div>

          <div className="info-row">
            <span className="label">Kontak :</span> <FaPhone />{" "}
            {formData.telepon || "Belum diisi"}
          </div>

          {/* Menampilkan Status Buka - SEKARANG FORMAT SIMPLE */}
          <div className="info-row">
            <span className="label">Jam Operasional :</span>{" "}
            <FaClock /> 
            <span className="operational-time">
              {statusBuka.hari && `${statusBuka.hari}, `}{statusBuka.waktu || "Loading..."}
            </span>
            <span
              style={{
                marginLeft: "10px",
                color: "white",
                background: statusBuka.warna,
                padding: "3px 8px",
                borderRadius: "10px",
                fontSize: "12px",
              }}
            >
              {statusBuka.status || "Loading..."}
            </span>
          </div>

          <div className="info-row">
            <span className="label">Fasilitas :</span> 
            {getFacilitiesText()}
          </div>
          <div className="info-row">
            <span className="label">Deskripsi :</span> 
            {formData.deskripsi || "Belum ada deskripsi"}
          </div>

          <div className="buttons">
            <button className="edit-btn" onClick={() => setShowEdit(true)}>
              {formData.nama ? 'Edit Profile ✏️' : 'Buat Profile 🆕'}
            </button>
            <button 
              className="debug-btn" 
              onClick={debugOperationalHours}
              title="Debug masalah jam operasional"
            >
              <FaBug /> Debug Jam
            </button>
            <button 
              className="refresh-btn" 
              onClick={() => {
                fetchCafeProfile();
                fetchGallery();
              }}
            >
              Refresh Data 🔄
            </button>
          </div>
        </div>
      </div>

      {/* Galeri */}
      <div className="gallery-section">
        <div className="gallery-card">
          <div className="gallery-header">
            <h3 className="gallery-title">Galeri Gambar</h3>
            <div className="gallery-counter">
              <span className={`counter-text ${maxGalleryReached ? 'counter-full' : ''}`}>
                {gallery.length}/10 gambar
                {maxGalleryReached && <FaExclamationTriangle className="warning-icon" />}
              </span>
            </div>
          </div>

          <div className="add-image">
            <label 
              htmlFor="fileInput" 
              className={`add-image-btn ${maxGalleryReached ? 'disabled' : ''}`}
              title={maxGalleryReached ? 'Galeri sudah penuh! Hapus beberapa gambar terlebih dahulu.' : 'Tambah gambar ke galeri'}
            >
              <FaPlus /> Tambah Gambar
            </label>
            <input
              id="fileInput"
              type="file"
              accept="image/*"
              style={{ display: "none" }}
              onChange={handleAddImage}
              disabled={maxGalleryReached}
            />
            {maxGalleryReached && (
              <div className="max-warning">
                <FaExclamationTriangle /> Galeri sudah penuh! Maksimal 10 gambar.
              </div>
            )}
          </div>

          <div className="gallery-grid">
            {gallery.length === 0 ? (
              <div className="empty-gallery">
                <p>Belum ada gambar di galeri</p>
                <small>Maksimal 10 gambar</small>
              </div>
            ) : (
              gallery.map((img, index) => (
                <div key={img.id} className="gallery-item">
                  <img 
                    src={getFullImageUrl(img.src)} 
                    alt={`Cafe ${index + 1}`} 
                    onError={(e) => {
                      e.target.src = "../../src/assets/LOGO.jpg";
                    }}
                  />
                  <button
                    className="delete-btn"
                    onClick={() => handleDeleteImage(img.id)}
                    title="Hapus gambar dari galeri"
                  >
                    <FaTrash />
                  </button>
                  <span className="pagination">{img.urutan}/{gallery.length}</span>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* Popup Edit */}
      {showEdit && (
        <div className="popup-overlay">
          <div className="popup-content">
            <button className="close-btn" onClick={() => setShowEdit(false)}>
              <FaTimes />
            </button>
            <h2>{formData.nama ? 'Edit Profile Cafe' : 'Buat Profile Cafe Baru'}</h2>

            <label>Nama Cafe:</label>
            <input name="nama" value={formData.nama} onChange={handleInputChange} placeholder="Masukkan nama cafe" />

            <label>Alamat:</label>
            <input name="alamat" value={formData.alamat} onChange={handleInputChange} placeholder="Masukkan alamat cafe" />

            <label>Nomor Telepon:</label>
            <input name="telepon" value={formData.telepon} onChange={handleInputChange} placeholder="Masukkan nomor telepon" />

            <label>Deskripsi:</label>
            <textarea
              name="deskripsi"
              value={formData.deskripsi}
              onChange={handleInputChange}
              rows="3"
              placeholder="Masukkan deskripsi cafe"
            />

            <label>Sosial Media:</label>
            {formData.social_media.map(sm => (
              <div key={sm.platform} className="social-input">
                <span className="platform-label">{sm.platform}:</span>
                <input
                  name={`social_${sm.platform}`}
                  value={sm.url}
                  onChange={handleInputChange}
                  placeholder={`https://${sm.platform}.com/username`}
                />
              </div>
            ))}

            <label>Jam Operasional:</label>
            {formData.operational_hours.map(hour => (
              <div key={hour.hari} className="day-time-row">
                <span className="day-label">
                  {hour.hari.charAt(0).toUpperCase() + hour.hari.slice(1)}:
                </span>
                <input
                  type="time"
                  name={`hour_${hour.hari}_buka`}
                  value={hour.buka || "08:00"}
                  onChange={handleInputChange}
                  required
                />
                <span> - </span>
                <input
                  type="time"
                  name={`hour_${hour.hari}_tutup`}
                  value={hour.tutup || "17:00"}
                  onChange={handleInputChange}
                  required
                />
              </div>
            ))}

            <label>Fasilitas:</label>
            <div className="facilities-grid">
              {formData.facilities.map(fac => (
                <div key={fac.nama_fasilitas} className="facility-item">
                  <input
                    type="checkbox"
                    name={`facility_${fac.nama_fasilitas}`}
                    checked={fac.tersedia}
                    onChange={handleInputChange}
                    id={`facility_${fac.nama_fasilitas}`}
                  />
                  <label htmlFor={`facility_${fac.nama_fasilitas}`}>
                    {fac.nama_fasilitas === 'wifi' ? 'Wifi' :
                     fac.nama_fasilitas === 'indoor' ? 'Indoor' :               
                     fac.nama_fasilitas === 'smoking_area' ? 'Smoking Area' :
                     fac.nama_fasilitas === 'indoor' ? 'Indoor' :
                     fac.nama_fasilitas === 'parkiran' ? 'Parkiran' :
                     fac.nama_fasilitas === 'musholla' ? 'Musholla' :
                     fac.nama_fasilitas === 'spot_foto' ? 'Spot Foto' :
                     fac.nama_fasilitas === 'live_music' ? 'Live Music' :
                     fac.nama_fasilitas === 'outdoor' ? 'Outdoor' :
                     fac.nama_fasilitas === 'air_conditioner' ? 'AC' : fac.nama_fasilitas}
                  </label>
                </div>
              ))}
            </div>

            <div className="popup-buttons">
              <button 
                className="cancel-btn" 
                onClick={() => setShowEdit(false)}
                disabled={isSaving}
              >
                Batal ❌
              </button>
              <button 
                className="save-btn" 
                onClick={handleSave}
                disabled={isSaving}
              >
                {isSaving ? 'Menyimpan...' : 'Simpan 💾'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default ProfileCafe;