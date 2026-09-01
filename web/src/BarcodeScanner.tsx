import { useEffect, useRef, useState } from "react";

interface DetectedCode { rawValue?: string }
type Detector = { detect: (source: HTMLVideoElement) => Promise<DetectedCode[]> };
type DetectorCtor = new (options?: { formats?: string[] }) => Detector;

type Props = { onDetected: (value: string) => void; onClose: () => void };

export function BarcodeScanner({ onDetected, onClose }: Props) {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const rafRef = useRef<number | null>(null);
  const [status, setStatus] = useState("Starting camera…");
  const [manual, setManual] = useState("");
  const [supported, setSupported] = useState(true);

  useEffect(() => {
    let active = true;
    const Detector = (globalThis as unknown as { BarcodeDetector?: DetectorCtor }).BarcodeDetector;
    if (!Detector) {
      setSupported(false);
      setStatus("Camera barcode detection is not available in this browser. Enter the code manually.");
      return () => { active = false; };
    }
    if (!navigator.mediaDevices?.getUserMedia) {
      setSupported(false);
      setStatus("Camera access is unavailable here. Enter the code manually.");
      return () => { active = false; };
    }

    const detector = new Detector({ formats: ["ean_13", "ean_8", "upc_a", "upc_e", "code_128", "code_39", "qr_code"] });
    void navigator.mediaDevices.getUserMedia({ video: { facingMode: { ideal: "environment" } }, audio: false })
      .then((stream) => {
        if (!active) { stream.getTracks().forEach((track) => track.stop()); return; }
        streamRef.current = stream;
        if (videoRef.current) {
          videoRef.current.srcObject = stream;
          void videoRef.current.play();
        }
        setStatus("Point the camera at a barcode or QR code.");

        const scan = async () => {
          if (!active || !videoRef.current) return;
          try {
            const results = await detector.detect(videoRef.current);
            const value = results.map((item) => item.rawValue?.trim()).find((item) => item);
            if (value) { onDetected(value); return; }
          } catch {
            setStatus("Camera is active. Move closer or improve the lighting.");
          }
          rafRef.current = window.requestAnimationFrame(() => { void scan(); });
        };
        rafRef.current = window.requestAnimationFrame(() => { void scan(); });
      })
      .catch(() => {
        if (active) { setSupported(false); setStatus("Camera permission was not granted. Enter the code manually."); }
      });

    return () => {
      active = false;
      if (rafRef.current !== null) window.cancelAnimationFrame(rafRef.current);
      streamRef.current?.getTracks().forEach((track) => track.stop());
      streamRef.current = null;
    };
  }, [onDetected]);

  function submitManual(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const value = manual.trim();
    if (value) onDetected(value);
  }

  return <div className="modal-backdrop" role="presentation">
    <section className="modal-card scanner-modal" role="dialog" aria-modal="true" aria-labelledby="barcode-scanner-title">
      <div className="panel-heading"><div><p className="eyebrow">Scanner</p><h2 id="barcode-scanner-title">Scan barcode or QR</h2></div><button className="secondary-button" type="button" onClick={onClose}>Close</button></div>
      <div className="scanner-stage">{supported ? <video ref={videoRef} className="scanner-video" playsInline muted aria-label="Barcode camera preview" /> : <div className="scanner-fallback"><strong>Manual barcode entry</strong><p>{status}</p></div>}</div>
      <p className="scanner-status" aria-live="polite">{status}</p>
      <form className="catalog-form scanner-manual" onSubmit={submitManual}><label><span>Barcode / QR value</span><input value={manual} onChange={(event) => setManual(event.target.value)} inputMode="numeric" autoComplete="off" maxLength={128} placeholder="Enter scanned value" /></label><div className="form-actions"><button className="primary-button" type="submit" disabled={!manual.trim()}>Use code</button></div></form>
    </section>
  </div>;
}
