# Muhasebe ve Is Takip Sistemi

Bu proje, isletmelerin finansal verilerini yonetmek, satis ve gider takibi yapmak ve isletme verimliligini analiz etmek amaciyla gelistirilmis bir web uygulamasidir. Ozellikle perde ve tekstil sektorundeki isletmelerin ihtiyaclarina yonelik ozellestirilmis moduller icermektedir.

## Proje Hakkinda

Uygulama, modern web teknolojileri kullanilarak full-stack bir mimari ile insa edilmistir. Isletme sahiplerinin gunluk islemlerini kolayca kayit altina almasini, gecmise donuk verileri raporlamasini ve dashboard uzerinden finansal durumu izlemesini saglar.

## Teknik Detaylar

### Frontend
- React: Kullanici arayuzu bilesenlerinin gelistirilmesi.
- Vite: Hizli gelistirme ortami ve build islemleri.
- Tailwind CSS: Modern ve responsive tasarim katmani.
- Recharts: Satis ve harcama verilerinin grafiksel gosterimi.
- Lucide React: Uygulama ici ikon bilesenleri.

### Backend
- Go (Golang): Yuksek performansli API sunucusu.
- Fiber v2: Hizli ve moduler web framework yapisi.
- Go Modules: Bagimlilik yonetimi.

## Klasor Yapisi

- perde-backend/: Isletme mantiginin ve finansal API uclarinin yonetildigi Go projesi.
- personel-takip/: Calisan verimliligi ve personel yonetimi icin ozellesmis moduller.
- src/: React frontend kaynak kodlari (Bilesenler, sayfalar ve servisler).
- public/: Statik varliklar ve dosyalar.

## Kurulum ve Calistirma

### Frontend Islemleri
1. Ana dizine gidin.
2. Bagimliliklari yukleyin:
   npm install
3. Uygulamayi gelistirme modunda baslatin:
   npm run dev

### Backend Islemleri
1. perde-backend dizinine gidin.
2. Gerekli Go paketlerini yukleyin:
   go mod tidy
3. Sunucuyu calistirin:
   go run main.go

## Not
Bu proje aktif gelistirme asamasindadir ve isletme gereksinimlerine gore ozellestirilmistir.
