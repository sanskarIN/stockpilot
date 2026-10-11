let stream = null;
let detector = null;
let animationFrame = null;
let active = false;

export async function startScanner(video, onDetected, onStatus) {
  stopScanner();
  const BarcodeDetectorCtor = globalThis.BarcodeDetector;
  if (!BarcodeDetectorCtor) throw new Error("Barcode scanning is not available in this browser. Enter the value manually.");
  if (!navigator.mediaDevices?.getUserMedia) throw new Error("Camera access is unavailable in this browser. Enter the value manually.");
  detector = new BarcodeDetectorCtor({ formats: ["qr_code", "ean_13", "ean_8", "upc_a", "upc_e", "code_39", "code_128"] });
  stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: { ideal: "environment" } }, audio: false });
  video.srcObject = stream;
  await video.play();
  active = true;
  onStatus("Point the camera at a barcode or QR code.");
  const scan = async () => {
    if (!active) return;
    try {
      const results = await detector.detect(video);
      const value = results.map((item) => item.rawValue?.trim()).find(Boolean);
      if (value) { onDetected(value); stopScanner(); return; }
    } catch { onStatus("Camera active. Move closer or improve lighting."); }
    animationFrame = requestAnimationFrame(() => void scan());
  };
  animationFrame = requestAnimationFrame(() => void scan());
}

export function stopScanner() {
  active = false;
  if (animationFrame !== null) cancelAnimationFrame(animationFrame);
  animationFrame = null;
  stream?.getTracks().forEach((track) => track.stop());
  stream = null;
  detector = null;
}
