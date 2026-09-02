import { useState, type ChangeEvent } from "react";
import { APIError, stockpilotAPI } from "./api";

type Props = { onClose: () => void; onSessionExpired: () => void; onImported: () => Promise<void> };
type Validation = Awaited<ReturnType<typeof stockpilotAPI.validateProductImport>>;

export function ProductImportPanel({ onClose, onSessionExpired, onImported }: Props) {
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<Validation | null>(null);
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState("");

  function chooseFile(event: ChangeEvent<HTMLInputElement>) {
    setFile(event.target.files?.[0] ?? null);
    setResult(null);
    setMessage("");
  }

  async function validate() {
    if (!file) { setMessage("Choose a CSV file first."); return; }
    if (!file.name.toLowerCase().endsWith(".csv")) { setMessage("Select a .csv file."); return; }
    setBusy(true); setMessage("");
    try {
      const validation = await stockpilotAPI.validateProductImport(file);
      setResult(validation);
    } catch (error) {
      if (error instanceof APIError && error.status === 401) { onSessionExpired(); return; }
      setMessage(error instanceof Error ? error.message : "Could not validate the CSV file.");
    } finally { setBusy(false); }
  }

  async function closeAfterReview() {
    if (result?.errorRows === 0) await onImported();
    onClose();
  }

  return <div className="modal-backdrop" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget && !busy) onClose(); }}>
    <section className="modal-card import-panel" role="dialog" aria-modal="true" aria-labelledby="product-import-title">
      <div className="panel-heading"><div><p className="eyebrow">Catalog import</p><h2 id="product-import-title">Validate products from CSV</h2></div><button className="secondary-button" type="button" onClick={onClose} disabled={busy}>Close</button></div>
      <p className="muted">Dry-run validation checks product fields, duplicate SKUs/barcodes, category and supplier references, and existing SKUs. Nothing is written to the database.</p>
      <label><span>CSV file</span><input type="file" accept=".csv,text/csv" onChange={chooseFile} disabled={busy} /></label>
      {file && <p className="table-secondary">Selected: {file.name} · {(file.size / 1024).toFixed(1)} KB</p>}
      {message && <p className="form-error" role="alert">{message}</p>}
      {result && <section className="import-result" aria-live="polite"><div className="import-summary"><strong>{result.validRows} valid rows</strong><strong>{result.errorRows} errors</strong></div>{result.errors.length > 0 && <div className="table-wrap"><table><thead><tr><th>Row</th><th>Problem</th></tr></thead><tbody>{result.errors.map((error, index) => <tr key={`${error.row}-${index}`}><td>{error.row}</td><td>{error.message}</td></tr>)}</tbody></table></div>}{result.errorRows === 0 && <p className="notice">Dry run passed. The file is ready for the write step in the next import milestone.</p>}</section>}
      <div className="form-actions"><button className="secondary-button" type="button" onClick={onClose} disabled={busy}>Cancel</button><button className="primary-button" type="button" onClick={result ? () => void closeAfterReview() : () => void validate()} disabled={busy || !file}>{busy ? "Validating…" : result ? "Done" : "Validate CSV"}</button></div>
    </section>
  </div>;
}
