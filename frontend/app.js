/* ============================================
   CuidaBien - Frontend JavaScript
   ============================================ */

const API_BASE = 'http://localhost:8098';

function asArray(value) {
  if (Array.isArray(value)) return value;
  if (value && Array.isArray(value.data)) return value.data;
  return [];
}

function listText(value) {
  if (Array.isArray(value)) return value.length ? value.join(', ') : '--';
  return value || '--';
}

function shortDate(value) {
  if (!value) return '--';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleDateString('es-ES');
}

function patientSelectOptions() {
  return [
    ['P001', 'Maria Garcia'],
    ['P002', 'Juan Lopez'],
    ['P003', 'Ana Martinez'],
  ].map(([id, name]) => `<option value="${id}">${name} (${id})</option>`).join('');
}

function csvList(value) {
  return value.split(',').map(v => v.trim()).filter(Boolean);
}

async function postJSON(path, payload) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || data.message || 'No se pudo guardar');
  return data;
}

// ============================================
// Initialization
// ============================================

document.addEventListener('DOMContentLoaded', () => {
  updateDate();
  checkApiStatus();
  loadStats();
});

function updateDate() {
  const now = new Date();
  const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
  document.getElementById('currentDate').textContent = now.toLocaleDateString('es-ES', options);
}

async function checkApiStatus() {
  try {
    const res = await fetch(`${API_BASE}/health`);
    const dot = document.querySelector('.status-dot');
    const text = document.querySelector('#apiStatus');
    if (res.ok) {
      dot.style.background = '#4ADE80';
      text.innerHTML = '<span class="status-dot"></span> Conectado';
    } else {
      dot.style.background = '#EF4444';
      text.innerHTML = '<span class="status-dot" style="background:#EF4444"></span> Error';
    }
  } catch {
    const text = document.querySelector('#apiStatus');
    text.innerHTML = '<span class="status-dot" style="background:#EF4444"></span> Desconectado';
  }
}

async function loadStats() {
  try {
    const [citasRes, alimRes, reportesRes, emergRes] = await Promise.all([
      fetch(`${API_BASE}/api/cita-medica`).then(r => r.json()).catch(() => []),
      fetch(`${API_BASE}/api/alimentacion/resumen`).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/api/reportes/resumen`).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/api/metrics`).then(r => r.json()).catch(() => null),
    ]);

    const hoy = new Date().toISOString().split('T')[0];
    const citas = asArray(citasRes);
    const citasHoy = citas.filter(c => c.fecha === hoy && (c.estado === 'pendiente' || c.estado === 'confirmada')).length;
    document.getElementById('statCitas').textContent = citasHoy;

    if (reportesRes && reportesRes.resumen_general) {
      document.getElementById('statPacientes').textContent = reportesRes.resumen_general.total_pacientes || 0;
      document.getElementById('statAlertas').textContent = reportesRes.resumen_general.total_alertas_pendientes || 0;
    }

    if (emergRes) {
      document.getElementById('statAlertas').textContent = emergRes.alertas_activas_hoy || 0;
    }

    if (alimRes) {
      document.getElementById('statComidas').textContent = `${alimRes.comidas_hechas || 0}/${alimRes.comidas_total || 3}`;
      document.getElementById('badgeAlim').textContent = `${alimRes.comidas_hechas || 0} de ${alimRes.comidas_total || 3}`;
    }
  } catch {
    document.getElementById('statCitas').textContent = '--';
  }
}

function openEmergencySelector() {
  const backdrop = document.getElementById('modalBackdrop');
  const modal = document.getElementById('modal');
  const title = document.getElementById('modalTitle');
  const body = document.getElementById('modalBody');

  title.textContent = 'Activar Emergencia';
  body.innerHTML = `<form class="inline-form danger-form" onsubmit="createQuickEmergency(event)">
    <h4>Seleccione el paciente</h4>
    <div class="form-grid">
      <label>Paciente<select name="paciente_id" required>${patientSelectOptions()}</select></label>
    </div>
    <button class="form-button danger-button" type="submit">Activar SOS</button>
  </form>`;
  backdrop.classList.add('active');
  modal.classList.add('active');
}

async function createQuickEmergency(event) {
  event.preventDefault();
  const { paciente_id } = Object.fromEntries(new FormData(event.target).entries());
  try {
    await postJSON('/api/alerts', {
      paciente_id,
      nivel: 'critico',
      mensaje: 'Boton de emergencia activado desde el panel CuidaBien',
    });
    await loadStats();
    openModal('emergencia');
  } catch (err) {
    alert(`No se pudo activar la emergencia: ${err.message}`);
  }
}

// ============================================
// Modal Management
// ============================================

function openModal(service) {
  const backdrop = document.getElementById('modalBackdrop');
  const modal = document.getElementById('modal');
  const title = document.getElementById('modalTitle');
  const body = document.getElementById('modalBody');

  backdrop.classList.add('active');
  modal.classList.add('active');

  const loaders = {
    'citas': loadCitasModal,
    'salud': loadSaludModal,
    'vitales': loadVitalesModal,
    'alimentacion': loadAlimentacionModal,
    'estado-animo': loadAnimoModal,
    'recordatorios': loadRecordatoriosModal,
    'ejercicios': loadEjerciciosModal,
    'actividad': loadActividadModal,
    'reportes': loadReportesModal,
    'emergencia': loadEmergenciaModal,
    'cuidadores': loadCuidadoresModal,
    'reportes-medicos': loadReportesMedicosModal,
  };

  title.textContent = getTitleForService(service);
  body.innerHTML = '<div class="loading-spinner"><div class="spinner"></div><p>Cargando datos...</p></div>';

  if (loaders[service]) {
    loaders[service](body);
  }
}

function closeModal() {
  document.getElementById('modalBackdrop').classList.remove('active');
  document.getElementById('modal').classList.remove('active');
}

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') closeModal();
});

function getTitleForService(service) {
  const titles = {
    'citas': 'Citas Medicas',
    'salud': 'Informacion de Salud',
    'vitales': 'Monitoreo de Signos Vitales',
    'alimentacion': 'Alimentacion e Hidratacion',
    'estado-animo': 'Estado de Animo',
    'recordatorios': 'Recordatorios de Medicamentos',
    'ejercicios': 'Estimulacion Cognitiva',
    'actividad': 'Actividad Fisica',
    'reportes': 'Reportes y Dashboard',
    'emergencia': 'Contacto de Emergencia',
    'cuidadores': 'Cuidadores',
    'reportes-medicos': 'Reportes Medicos',
  };
  return titles[service] || service;
}

function badge(estado) {
  const cls = {
    'completada': 'badge-success', 'cumplida': 'badge-success', 'activo': 'badge-success',
    'confirmada': 'badge-info', 'pendiente': 'badge-warning',
    'cancelada': 'badge-error', 'omitida': 'badge-error', 'suspendido': 'badge-error',
  };
  return `<span class="badge ${cls[estado] || 'badge-info'}">${estado}</span>`;
}

// ============================================
// CITAS
// ============================================

async function loadCitasModal(body) {
  try {
    const citasResponse = await fetch(`${API_BASE}/api/cita-medica`).then(r => r.json());
    const citas = asArray(citasResponse);
    const pacientes = await fetch(`${API_BASE}/api/cita-medica/pacientes`).then(r => r.json()).catch(() => []);
    const doctores = await fetch(`${API_BASE}/api/cita-medica/doctores`).then(r => r.json()).catch(() => []);

    let html = '';

    if (doctores.length > 0) {
      html += `<div class="detail-section"><h4>Doctores Disponibles</h4><div class="detail-grid">`;
      doctores.forEach(d => {
        html += `<div class="detail-item"><div class="label">${d.especialidad}</div><div class="value">${d.nombre}</div></div>`;
      });
      html += `</div></div>`;
    }

    html += `<form class="inline-form" onsubmit="createCita(event)">`;
    html += `<h4>Agregar Cita</h4><div class="form-grid">`;
    html += `<label>Paciente<select name="paciente_id" required>${pacientes.map(p => `<option value="${p.id}">${p.nombre} (${p.id})</option>`).join('')}</select></label>`;
    html += `<label>Doctor<select name="doctor_id" required>${doctores.map(d => `<option value="${d.id}">${d.nombre} - ${d.especialidad}</option>`).join('')}</select></label>`;
    html += `<label>Fecha<input name="fecha" type="date" required></label>`;
    html += `<label>Hora<input name="hora" type="time" required></label>`;
    html += `<label>Prioridad<select name="prioridad"><option value="normal">Normal</option><option value="control">Control</option><option value="urgente">Urgente</option></select></label>`;
    html += `<label>Motivo<input name="motivo" placeholder="Control, dolor, chequeo..."></label>`;
    html += `</div><button class="form-button" type="submit">Guardar cita</button></form>`;

    html += `<div class="detail-section"><h4>Citas Registradas (${Array.isArray(citas) ? citas.length : 0})</h4>`;
    if (Array.isArray(citas) && citas.length > 0) {
      html += `<table class="data-table"><thead><tr><th>ID</th><th>Paciente</th><th>Doctor</th><th>Fecha</th><th>Hora</th><th>Estado</th><th>Prioridad</th></tr></thead><tbody>`;
      citas.forEach(c => {
        const paciente = pacientes.find(p => p.id === c.paciente_id);
        const doctor = doctores.find(d => d.id === c.doctor_id);
        html += `<tr><td>${c.id}</td><td>${paciente ? paciente.nombre : c.paciente_id}</td><td>${doctor ? doctor.nombre : c.doctor_id}</td><td>${c.fecha}</td><td>${c.hora}</td><td>${badge(c.estado)}</td><td>${badge(c.prioridad)}</td></tr>`;
      });
      html += `</tbody></table>`;
    } else {
      html += `<div class="empty-state"><p>No hay citas registradas</p></div>`;
    }
    html += `</div>`;

    document.getElementById('badgeCitas').textContent = Array.isArray(citas) ? citas.length : 0;
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error al cargar citas: ${err.message}</p></div>`;
  }
}

async function createCita(event) {
  event.preventDefault();
  const form = event.target;
  const data = Object.fromEntries(new FormData(form).entries());
  try {
    await postJSON('/api/cita-medica', data);
    await loadStats();
    openModal('citas');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// INFORMACION DE SALUD
// ============================================

async function loadSaludModal(body) {
  try {
    const data = await fetch(`${API_BASE}/api/informacion-salud`).then(r => r.json());
    let html = '';

    html += `<form class="inline-form" onsubmit="createSalud(event)">`;
    html += `<h4>Agregar/actualizar ficha de salud</h4><div class="form-grid">`;
    html += `<label>ID Paciente<input name="paciente_id" placeholder="P001" required></label>`;
    html += `<label>Nombre Paciente<input name="nombre_paciente" placeholder="Maria Garcia" required></label>`;
    html += `<label>Diagnosticos<input name="diagnosticos" placeholder="hipertension, artritis"></label>`;
    html += `<label>Alergias<input name="alergias" placeholder="penicilina, mariscos"></label>`;
    html += `<label>Enfermedades Cronicas<input name="enfermedades_cronicas" placeholder="diabetes tipo 2"></label>`;
    html += `<label>Antecedentes Medicos<input name="antecedentes_medicos" placeholder="cirugia, marcapasos"></label>`;
    html += `</div><button class="form-button" type="submit">Guardar ficha</button></form>`;

    if (Array.isArray(data) && data.length > 0) {
      data.forEach(p => {
        html += `<div class="summary-card"><h5>${p.nombre_paciente || p.nombre || 'Paciente'}</h5>`;
        html += `<div class="detail-grid">`;
        const fields = [
          ['Diagnosticos', listText(p.diagnosticos)],
          ['Alergias', listText(p.alergias)],
          ['Enfermedades Cronicas', listText(p.enfermedades_cronicas)],
          ['Antecedentes Medicos', listText(p.antecedentes_medicos)],
        ];
        fields.forEach(([label, value]) => {
          html += `<div class="detail-item"><div class="label">${label}</div><div class="value">${value || '--'}</div></div>`;
        });
        html += `</div></div>`;
      });
    } else {
      html = `<div class="empty-state"><p>No hay registros de salud</p></div>`;
    }
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createSalud(event) {
  event.preventDefault();
  const form = event.target;
  const values = Object.fromEntries(new FormData(form).entries());
  const payload = {
    paciente_id: values.paciente_id,
    nombre_paciente: values.nombre_paciente,
    diagnosticos: csvList(values.diagnosticos || ''),
    alergias: csvList(values.alergias || ''),
    enfermedades_cronicas: csvList(values.enfermedades_cronicas || ''),
    antecedentes_medicos: csvList(values.antecedentes_medicos || ''),
  };
  try {
    await postJSON('/api/informacion-salud', payload);
    openModal('salud');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// SIGNOS VITALES
// ============================================

async function loadVitalesModal(body) {
  try {
    const ids = ['P001', 'P002', 'P003'];
    let html = `<form class="inline-form" onsubmit="createVitales(event)">`;
    html += `<h4>Agregar Signos Vitales</h4><div class="form-grid">`;
    html += `<label>Paciente<select name="id_adulto_mayor" required>${ids.map(id => `<option value="${id}">${id}</option>`).join('')}</select></label>`;
    html += `<label>Registrado por<input name="registrado_por" placeholder="Cuidador Ana"></label>`;
    html += `<label>Sistolica<input name="presion_sistolica" type="number" required></label>`;
    html += `<label>Diastolica<input name="presion_diastolica" type="number" required></label>`;
    html += `<label>Frecuencia<input name="frecuencia_cardiaca" type="number"></label>`;
    html += `<label>Temperatura<input name="temperatura" type="number" step="0.1"></label>`;
    html += `<label>Oxigenacion<input name="saturacion_oxigeno" type="number"></label>`;
    html += `<label>Observaciones<input name="observaciones" placeholder="Control de rutina"></label>`;
    html += `</div><button class="form-button" type="submit">Guardar signos</button></form>`;

    for (const id of ids) {
      const data = await fetch(`${API_BASE}/api/vitales/${id}/`).then(r => r.json()).catch(() => []);
      const latest = await fetch(`${API_BASE}/api/vitales/${id}/ultimo`).then(r => r.json()).catch(() => null);

      html += `<div class="detail-section"><h4>Adulto Mayor #${id}</h4>`;

      if (latest) {
        html += `<div class="detail-grid">`;
        const pa = latest.presion_sistolica && latest.presion_diastolica ? `${latest.presion_sistolica}/${latest.presion_diastolica} mmHg` : '--';
        html += `<div class="detail-item"><div class="label">Presion Arterial</div><div class="value">${pa}</div></div>`;
        html += `<div class="detail-item"><div class="label">Frecuencia Cardiaca</div><div class="value">${latest.frecuencia_cardiaca || '--'} lpm</div></div>`;
        html += `<div class="detail-item"><div class="label">Temperatura</div><div class="value">${latest.temperatura || '--'} C</div></div>`;
        html += `<div class="detail-item"><div class="label">Oxigenacion</div><div class="value">${latest.saturacion_oxigeno || '--'}%</div></div>`;
        html += `</div>`;
      }

      if (Array.isArray(data) && data.length > 0) {
        html += `<table class="data-table"><thead><tr><th>Fecha</th><th>Presion</th><th>Pulso</th><th>Temp</th><th>Estado</th></tr></thead><tbody>`;
        data.slice(0, 5).forEach(r => {
          const pa = r.presion_sistolica && r.presion_diastolica ? `${r.presion_sistolica}/${r.presion_diastolica}` : '--';
          html += `<tr><td>${shortDate(r.fecha_registro)}</td><td>${pa}</td><td>${r.frecuencia_cardiaca || '--'}</td><td>${r.temperatura || '--'}</td><td>${badge((r.evaluacion && r.evaluacion.estado_general) || 'normal')}</td></tr>`;
        });
        html += `</tbody></table>`;
      }
      html += `</div>`;
    }
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createVitales(event) {
  event.preventDefault();
  const values = Object.fromEntries(new FormData(event.target).entries());
  const payload = {
    id_adulto_mayor: values.id_adulto_mayor,
    registrado_por: values.registrado_por,
    presion_sistolica: Number(values.presion_sistolica),
    presion_diastolica: Number(values.presion_diastolica),
    observaciones: values.observaciones,
  };
  ['frecuencia_cardiaca', 'temperatura', 'saturacion_oxigeno'].forEach(key => {
    if (values[key] !== '') payload[key] = Number(values[key]);
  });
  try {
    await postJSON('/api/vitales', payload);
    openModal('vitales');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// ALIMENTACION
// ============================================

async function loadAlimentacionModal(body) {
  try {
    const [resumen, comidas, hidratacion, restricciones] = await Promise.all([
      fetch(`${API_BASE}/api/alimentacion/resumen`).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/api/alimentacion`).then(r => r.json()).catch(() => []),
      fetch(`${API_BASE}/api/alimentacion/hidratacion`).then(r => r.json()).catch(() => []),
      fetch(`${API_BASE}/api/alimentacion/restricciones`).then(r => r.json()).catch(() => []),
    ]);

    let html = '';

    html += `<form class="inline-form" onsubmit="createComida(event)">`;
    html += `<h4>Registrar Comida</h4><div class="form-grid">`;
    html += `<label>Tipo<select name="tipo_comida" required><option value="desayuno">Desayuno</option><option value="almuerzo">Almuerzo</option><option value="cena">Cena</option><option value="merienda">Merienda</option></select></label>`;
    html += `<label>Descripcion<input name="descripcion" placeholder="Sopa, pollo, frutas"></label>`;
    html += `</div><button class="form-button" type="submit">Guardar comida</button></form>`;

    html += `<form class="inline-form compact-form" onsubmit="createHidratacion(event)">`;
    html += `<h4>Registrar Hidratacion</h4><div class="form-grid">`;
    html += `<label>Cantidad<input name="cantidad" placeholder="1 vaso de agua" required></label>`;
    html += `</div><button class="form-button" type="submit">Guardar hidratacion</button></form>`;

    html += `<form class="inline-form compact-form" onsubmit="createRestriccion(event)">`;
    html += `<h4>Agregar Restriccion</h4><div class="form-grid">`;
    html += `<label>Descripcion<input name="descripcion" placeholder="Sin sal, diabetico, alergia a lacteos" required></label>`;
    html += `</div><button class="form-button" type="submit">Guardar restriccion</button></form>`;

    if (resumen) {
      const nivelColor = resumen.nivel_alerta === 'ok' ? 'badge-success' : resumen.nivel_alerta === 'atencion' ? 'badge-warning' : 'badge-error';
      html += `<div class="summary-card"><h5>Resumen del Dia</h5><div class="detail-grid">`;
      html += `<div class="detail-item"><div class="label">Comidas Realizadas</div><div class="value">${resumen.comidas_hechas} / ${resumen.comidas_total}</div></div>`;
      html += `<div class="detail-item"><div class="label">Nivel de Alerta</div><div class="value"><span class="badge ${nivelColor}">${resumen.nivel_alerta}</span></div></div>`;
      html += `<div class="detail-item"><div class="label">Comidas Saltadas</div><div class="value">${resumen.comidas ? resumen.comidas.filter(c => c.saltada).length : 0}</div></div>`;
      html += `</div></div>`;
    }

    if (Array.isArray(restricciones) && restricciones.length > 0) {
      html += `<div class="detail-section"><h4>Restricciones Alimentarias</h4><div class="detail-grid">`;
      restricciones.forEach(r => {
        html += `<div class="detail-item"><div class="label">${r.id}</div><div class="value">${r.descripcion}</div></div>`;
      });
      html += `</div></div>`;
    }

    if (Array.isArray(comidas) && comidas.length > 0) {
      html += `<div class="detail-section"><h4>Comidas de Hoy</h4><table class="data-table"><thead><tr><th>ID</th><th>Tipo</th><th>Descripcion</th><th>Hora</th></tr></thead><tbody>`;
      comidas.forEach(c => {
        const hora = c.hora ? new Date(c.hora).toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' }) : '--';
        html += `<tr><td>${c.id}</td><td>${c.tipo_comida}</td><td>${c.descripcion || '--'}</td><td>${hora}</td></tr>`;
      });
      html += `</tbody></table></div>`;
    }

    if (Array.isArray(hidratacion) && hidratacion.length > 0) {
      html += `<div class="detail-section"><h4>Hidratacion</h4><table class="data-table"><thead><tr><th>Hora</th><th>Cantidad</th></tr></thead><tbody>`;
      hidratacion.forEach(h => {
        const hora = h.hora ? new Date(h.hora).toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' }) : '--';
        html += `<tr><td>${hora}</td><td>${h.cantidad || '--'}</td></tr>`;
      });
      html += `</tbody></table></div>`;
    }

    body.innerHTML = html || `<div class="empty-state"><p>No hay datos de alimentacion</p></div>`;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createComida(event) {
  event.preventDefault();
  const data = Object.fromEntries(new FormData(event.target).entries());
  try {
    await postJSON('/api/alimentacion', data);
    await loadStats();
    openModal('alimentacion');
  } catch (err) {
    alert(err.message);
  }
}

async function createHidratacion(event) {
  event.preventDefault();
  const data = Object.fromEntries(new FormData(event.target).entries());
  try {
    await postJSON('/api/alimentacion/hidratacion', data);
    openModal('alimentacion');
  } catch (err) {
    alert(err.message);
  }
}

async function createRestriccion(event) {
  event.preventDefault();
  const data = Object.fromEntries(new FormData(event.target).entries());
  try {
    await postJSON('/api/alimentacion/restricciones', data);
    openModal('alimentacion');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// ESTADO DE ANIMO
// ============================================

async function loadAnimoModal(body) {
  try {
    const data = await fetch(`${API_BASE}/api/estado-animo`).then(r => r.json()).catch(() => []);
    let html = '';

    html += `<form class="inline-form" onsubmit="createAnimo(event)">`;
    html += `<h4>Registrar Estado de Animo</h4><div class="form-grid">`;
    html += `<label>Nivel<select name="nivel" required><option value="5">5 - Feliz</option><option value="4">4 - Tranquilo</option><option value="3">3 - Neutral</option><option value="2">2 - Triste</option><option value="1">1 - Muy triste</option></select></label>`;
    html += `<label>Emocion<input name="emocion" placeholder="Feliz, tranquilo, preocupado" required></label>`;
    html += `<label>Comentario<input name="comentario" placeholder="Como se sintio hoy"></label>`;
    html += `<label>Fecha<input name="fecha" type="date"></label>`;
    html += `</div><button class="form-button" type="submit">Guardar estado</button></form>`;

    if (Array.isArray(data) && data.length > 0) {
      html += `<div class="detail-section"><h4>Registros de Animo (${data.length})</h4>`;
      html += `<div class="detail-grid">`;
      data.forEach(r => {
        const emoji = r.nivel >= 4 ? '&#128522;' : r.nivel >= 3 ? '&#128528;' : '&#128533;';
        html += `<div class="detail-item"><div class="label">${r.fecha || r.timestamp || 'Hoy'}</div><div class="value">${emoji} ${r.emocion || r.estado || '--'} (Nivel ${r.nivel})</div></div>`;
      });
      html += `</div></div>`;

      document.getElementById('badgeAnimo').textContent = data.length + ' registros';
    } else {
      html = `<div class="empty-state"><p>No hay registros de animo</p></div>`;
    }
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createAnimo(event) {
  event.preventDefault();
  const values = Object.fromEntries(new FormData(event.target).entries());
  const payload = { nivel: Number(values.nivel), emocion: values.emocion, comentario: values.comentario };
  if (values.fecha) payload.fecha = values.fecha;
  try {
    await postJSON('/api/estado-animo', payload);
    openModal('estado-animo');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// RECORDATORIOS
// ============================================

async function loadRecordatoriosModal(body) {
  try {
    const data = await fetch(`${API_BASE}/api/recordatorios-medicamentos`).then(r => r.json()).catch(() => []);
    let html = '';

    html += `<form class="inline-form" onsubmit="createRecordatorio(event)">`;
    html += `<h4>Agregar Recordatorio</h4><div class="form-grid">`;
    html += `<label>ID Paciente<input name="adulto_mayor_id" placeholder="P001" required></label>`;
    html += `<label>Paciente<input name="nombre_paciente" placeholder="Maria Garcia" required></label>`;
    html += `<label>Medicamento<input name="medicamento" placeholder="Losartan" required></label>`;
    html += `<label>Dosis<input name="dosis" placeholder="1 tableta" required></label>`;
    html += `<label>Hora<input name="hora" type="time" required></label>`;
    html += `<label>Frecuencia<input name="frecuencia" placeholder="diaria" required></label>`;
    html += `</div><button class="form-button" type="submit">Guardar recordatorio</button></form>`;

    if (Array.isArray(data) && data.length > 0) {
      html += `<div class="detail-section"><h4>Recordatorios Activos (${data.length})</h4>`;
      html += `<table class="data-table"><thead><tr><th>ID</th><th>Paciente</th><th>Medicamento</th><th>Hora</th><th>Frecuencia</th><th>Activo</th></tr></thead><tbody>`;
      data.forEach(r => {
        const activo = r.activo !== false;
        html += `<tr><td>${r.id}</td><td>${r.nombre_paciente || r.paciente_id || r.paciente_nombre || '--'}</td><td>${r.medicamento || r.nombre || '--'}</td><td>${r.hora || '--'}</td><td>${r.frecuencia || '--'}</td><td>${activo ? badge('activo') : badge('suspendido')}</td></tr>`;
      });
      html += `</tbody></table></div>`;
      document.getElementById('badgeRec').textContent = data.length;
    } else {
      html = `<div class="empty-state"><p>No hay recordatorios registrados</p></div>`;
    }
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createRecordatorio(event) {
  event.preventDefault();
  const data = Object.fromEntries(new FormData(event.target).entries());
  try {
    await postJSON('/api/recordatorios-medicamentos', data);
    openModal('recordatorios');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// ESTIMULACION COGNITIVA
// ============================================

async function loadEjerciciosModal(body) {
  try {
    const [ejercicios, resumen] = await Promise.all([
      fetch(`${API_BASE}/api/ejercicios`).then(r => r.json()).catch(() => []),
      fetch(`${API_BASE}/api/ejercicios/resumen`).then(r => r.json()).catch(() => null),
    ]);

    let html = '';

    html += `<form class="inline-form" onsubmit="createEjercicio(event)">`;
    html += `<h4>Registrar Ejercicio Cognitivo</h4><div class="form-grid">`;
    html += `<label>Tipo<select name="tipo" required><option value="memoria">Memoria</option><option value="trivia">Trivia</option><option value="sopa_letras">Sopa de letras</option><option value="calculo">Calculo</option></select></label>`;
    html += `</div><button class="form-button" type="submit">Guardar ejercicio</button></form>`;

    if (resumen) {
      const nivelActividad = resumen.hay_alerta ? 'requiere atencion' : 'al dia';
      html += `<div class="summary-card"><h5>Resumen de Actividad Cognitiva</h5><div class="detail-grid">`;
      html += `<div class="detail-item"><div class="label">Total Ejercicios</div><div class="value">${resumen.total_ejercicios || resumen.total || 0}</div></div>`;
      html += `<div class="detail-item"><div class="label">Nivel de Actividad</div><div class="value">${resumen.nivel_actividad || resumen.nivel || nivelActividad}</div></div>`;
      html += `<div class="detail-item"><div class="label">Ejercicios Hoy</div><div class="value">${resumen.ejercicios_hoy || 0}</div></div>`;
      html += `<div class="detail-item"><div class="label">Seguimiento</div><div class="value">${resumen.mensaje || '--'}</div></div>`;
      html += `</div></div>`;
    }

    if (Array.isArray(ejercicios) && ejercicios.length > 0) {
      html += `<div class="detail-section"><h4>Ejercicios Realizados</h4>`;
      html += `<table class="data-table"><thead><tr><th>Tipo</th><th>Fecha</th><th>Registro</th></tr></thead><tbody>`;
      ejercicios.forEach(e => {
        html += `<tr><td>${e.tipo || e.nombre || '--'}</td><td>${shortDate(e.fecha || e.fecha_hora)}</td><td>Completado</td></tr>`;
      });
      html += `</tbody></table></div>`;
      document.getElementById('badgeEjerc').textContent = ejercicios.length + ' ejercicios';
    }

    body.innerHTML = html || `<div class="empty-state"><p>No hay ejercicios registrados</p></div>`;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createEjercicio(event) {
  event.preventDefault();
  const data = Object.fromEntries(new FormData(event.target).entries());
  try {
    await postJSON('/api/ejercicios', data);
    openModal('ejercicios');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// ACTIVIDAD FISICA
// ============================================

async function loadActividadModal(body) {
  try {
    const data = await fetch(`${API_BASE}/api/actividad-fisica`).then(r => r.json()).catch(() => []);
    let html = '';

    html += `<form class="inline-form" onsubmit="createActividad(event)">`;
    html += `<h4>Registrar Actividad Fisica</h4><div class="form-grid">`;
    html += `<label>Paciente<input name="nombre_paciente" placeholder="Maria Garcia" required></label>`;
    html += `<label>Actividad<input name="tipo_actividad" placeholder="Caminata" required></label>`;
    html += `<label>Duracion (min)<input name="duracion_minutos" type="number" min="1" required></label>`;
    html += `<label>Intensidad<select name="intensidad" required><option value="baja">Baja</option><option value="moderada">Moderada</option><option value="alta">Alta</option></select></label>`;
    html += `<label>Fecha<input name="fecha" type="date" required></label>`;
    html += `<label>Estado<select name="estado" required><option value="completada">Completada</option><option value="pendiente">Pendiente</option><option value="cancelada">Cancelada</option></select></label>`;
    html += `<label>Observaciones<input name="observaciones" placeholder="Sin molestias"></label>`;
    html += `</div><button class="form-button" type="submit">Guardar actividad</button></form>`;

    if (Array.isArray(data) && data.length > 0) {
      html += `<div class="detail-section"><h4>Actividades Registradas (${data.length})</h4>`;
      html += `<table class="data-table"><thead><tr><th>ID</th><th>Actividad</th><th>Duracion</th><th>Intensidad</th><th>Fecha</th></tr></thead><tbody>`;
      data.forEach(a => {
        const duracion = a.duracion_minutos ? `${a.duracion_minutos} min` : (a.duracion || '--');
        html += `<tr><td>${a.id}</td><td>${a.tipo_actividad || a.actividad || a.tipo || a.nombre || '--'}</td><td>${duracion}</td><td>${a.intensidad || '--'}</td><td>${a.fecha || a.fecha_hora || '--'}</td></tr>`;
      });
      html += `</tbody></table></div>`;
      document.getElementById('badgeActiv').textContent = data.length + ' registros';
    } else {
      html = `<div class="empty-state"><p>No hay actividades registradas</p></div>`;
    }
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createActividad(event) {
  event.preventDefault();
  const values = Object.fromEntries(new FormData(event.target).entries());
  const payload = { ...values, duracion_minutos: Number(values.duracion_minutos) };
  try {
    await postJSON('/api/actividad-fisica', payload);
    openModal('actividad');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// REPORTES
// ============================================

async function loadReportesModal(body) {
  try {
    const dashboard = await fetch(`${API_BASE}/api/reportes/resumen`).then(r => r.json()).catch(() => null);

    let html = '';

    if (dashboard && dashboard.resumen_general) {
      const rg = dashboard.resumen_general;
      html += `<div class="summary-card"><h5>Dashboard General</h5><div class="detail-grid">`;
      html += `<div class="detail-item"><div class="label">Total Pacientes</div><div class="value">${rg.total_pacientes}</div></div>`;
      html += `<div class="detail-item"><div class="label">Citas Hoy</div><div class="value">${rg.total_citas_hoy}</div></div>`;
      html += `<div class="detail-item"><div class="label">Medicamentos Activos</div><div class="value">${rg.total_medicamentos_activos}</div></div>`;
      html += `<div class="detail-item"><div class="label">Alertas Pendientes</div><div class="value">${rg.total_alertas_pendientes}</div></div>`;
      html += `<div class="detail-item"><div class="label">Promedio Adherencia</div><div class="value">${rg.promedio_adherencia}%</div></div>`;
      html += `<div class="detail-item"><div class="label">Pacientes con Alertas</div><div class="value">${rg.pacientes_con_alertas}</div></div>`;
      html += `</div></div>`;
    }

    if (dashboard && dashboard.pacientes && dashboard.pacientes.length > 0) {
      html += `<div class="detail-section"><h4>Resumen por Paciente</h4>`;
      html += `<table class="data-table"><thead><tr><th>ID</th><th>Nombre</th><th>Citas Proximas</th><th>Meds Activos</th><th>Adherencia</th><th>Estado</th></tr></thead><tbody>`;
      dashboard.pacientes.forEach(p => {
        html += `<tr><td>${p.id}</td><td>${p.nombre}</td><td>${p.citas_proximas}</td><td>${p.medicamentos_activos}</td><td>${p.adherencia}%</td><td>${badge(p.estado)}</td></tr>`;
      });
      html += `</tbody></table></div>`;
    }

    body.innerHTML = html || `<div class="empty-state"><p>No hay datos de reportes</p></div>`;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

// ============================================
// CONTACTO DE EMERGENCIA
// ============================================

async function loadEmergenciaModal(body) {
  try {
    const contacts = await fetch(`${API_BASE}/api/contacts`).then(r => r.json()).catch(() => []);
    const alerts = await fetch(`${API_BASE}/api/alerts`).then(r => r.json()).catch(() => []);

    let html = '';

    html += `<form class="inline-form" onsubmit="createContacto(event)">`;
    html += `<h4>Agregar Contacto de Emergencia</h4><div class="form-grid">`;
    html += `<label>ID Paciente<input name="paciente_id" placeholder="P001" required></label>`;
    html += `<label>Nombre<input name="nombre" required></label>`;
    html += `<label>Telefono<input name="telefono" required></label>`;
    html += `<label>Relacion<input name="parentesco" placeholder="Hija, esposo, cuidador" required></label>`;
    html += `<label>Prioridad<input name="prioridad" type="number" min="1" value="1"></label>`;
    html += `<label>Principal<select name="principal"><option value="true">Si</option><option value="false">No</option></select></label>`;
    html += `</div><button class="form-button" type="submit">Guardar contacto</button></form>`;

    html += `<form class="inline-form compact-form danger-form" onsubmit="createAlerta(event)">`;
    html += `<h4>Activar Emergencia</h4><div class="form-grid">`;
    html += `<label>Paciente<select name="paciente_id" required>${patientSelectOptions()}</select></label>`;
    html += `</div><button class="form-button danger-button" type="submit">Activar SOS</button></form>`;

    if (Array.isArray(contacts) && contacts.length > 0) {
      html += `<div class="detail-section"><h4>Contactos de Emergencia (${contacts.length})</h4>`;
      html += `<table class="data-table"><thead><tr><th>Nombre</th><th>Relacion</th><th>Telefono</th><th>Paciente</th></tr></thead><tbody>`;
      contacts.forEach(c => {
        html += `<tr><td>${c.nombre || c.name || '--'}</td><td>${c.parentesco || c.relacion || c.relationship || '--'}</td><td>${c.telefono || c.phone || '--'}</td><td>${c.paciente_id || c.pacienteID || '--'}</td></tr>`;
      });
      html += `</tbody></table></div>`;
      document.getElementById('badgeEmerg').textContent = contacts.length + ' contactos';
    }

    if (Array.isArray(alerts) && alerts.length > 0) {
      html += `<div class="detail-section"><h4>Alertas de Emergencia</h4>`;
      html += `<table class="data-table"><thead><tr><th>ID</th><th>Tipo</th><th>Mensaje</th><th>Estado</th></tr></thead><tbody>`;
      alerts.forEach(a => {
        html += `<tr><td>${a.id}</td><td>${a.tipo || a.type || '--'}</td><td>${a.mensaje || a.message || '--'}</td><td>${badge(a.estado || a.status || 'activa')}</td></tr>`;
      });
      html += `</tbody></table></div>`;
    }

    body.innerHTML = html || `<div class="empty-state"><p>No hay contactos ni alertas de emergencia</p></div>`;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createContacto(event) {
  event.preventDefault();
  const values = Object.fromEntries(new FormData(event.target).entries());
  const payload = {
    paciente_id: values.paciente_id,
    nombre: values.nombre,
    telefono: values.telefono,
    parentesco: values.parentesco,
    prioridad: Number(values.prioridad || 1),
    principal: values.principal === 'true',
  };
  try {
    await postJSON('/api/contacts', payload);
    await loadStats();
    openModal('emergencia');
  } catch (err) {
    alert(err.message);
  }
}

async function createAlerta(event) {
  event.preventDefault();
  const { paciente_id } = Object.fromEntries(new FormData(event.target).entries());
  try {
    await postJSON('/api/alerts', {
      paciente_id,
      nivel: 'critico',
      mensaje: 'Emergencia activada desde el panel CuidaBien',
    });
    await loadStats();
    openModal('emergencia');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// CUIDADORES
// ============================================

async function loadCuidadoresModal(body) {
  try {
    const data = await fetch(`${API_BASE}/api/cuidadores`).then(r => r.json()).catch(() => []);
    let html = '';

    html += `<form class="inline-form" onsubmit="createCuidador(event)">`;
    html += `<h4>Agregar Cuidador</h4><div class="form-grid">`;
    html += `<label>Nombre<input name="nombre" required></label>`;
    html += `<label>Telefono<input name="telefono" required></label>`;
    html += `<label>Email<input name="email" type="email"></label>`;
    html += `<label>Relacion<input name="relacion" placeholder="Familiar, enfermero"></label>`;
    html += `<label>Horario<input name="horario_disponible" placeholder="Lunes a viernes"></label>`;
    html += `<label>Pacientes<input name="pacientes" placeholder="P001, P002"></label>`;
    html += `<label>Responsabilidad<select name="nivel_responsabilidad"><option value="alta">Alta</option><option value="media">Media</option><option value="apoyo">Apoyo</option></select></label>`;
    html += `</div><button class="form-button" type="submit">Guardar cuidador</button></form>`;

    if (Array.isArray(data) && data.length > 0) {
      html += `<div class="detail-section"><h4>Cuidadores Registrados (${data.length})</h4>`;
      html += `<table class="data-table"><thead><tr><th>ID</th><th>Nombre</th><th>Telefono</th><th>Pacientes</th><th>Horario</th><th>Responsabilidad</th></tr></thead><tbody>`;
      data.forEach(c => {
        html += `<tr><td>${c.id}</td><td>${c.nombre || c.name || '--'}</td><td>${c.telefono || c.phone || '--'}</td><td>${listText(c.pacientes)}</td><td>${c.horario_disponible || c.turno || c.horario || '--'}</td><td>${c.nivel_responsabilidad || '--'}</td></tr>`;
      });
      html += `</tbody></table></div>`;
      document.getElementById('badgeCuid').textContent = data.length + ' cuidadores';
    } else {
      html = `<div class="empty-state"><p>No hay cuidadores registrados</p></div>`;
    }
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}

async function createCuidador(event) {
  event.preventDefault();
  const values = Object.fromEntries(new FormData(event.target).entries());
  const payload = {
    nombre: values.nombre,
    telefono: values.telefono,
    email: values.email,
    relacion: values.relacion,
    horario_disponible: values.horario_disponible,
    pacientes: csvList(values.pacientes || ''),
    nivel_responsabilidad: values.nivel_responsabilidad,
  };
  try {
    await postJSON('/api/cuidadores', payload);
    openModal('cuidadores');
  } catch (err) {
    alert(err.message);
  }
}

// ============================================
// REPORTES MEDICOS
// ============================================

async function loadReportesMedicosModal(body) {
  try {
    const data = await fetch(`${API_BASE}/api/reportes-medicos/semanal`).then(r => r.json()).catch(() => []);
    let html = '';

    if (Array.isArray(data) && data.length > 0) {
      data.forEach(r => {
        html += `<div class="summary-card"><h5>${r.nombre || r.paciente_nombre || 'Paciente'} - ${r.periodo || 'Periodo actual'}</h5>`;
        html += `<div class="detail-grid">`;
        html += `<div class="detail-item"><div class="label">Citas Programadas</div><div class="value">${r.citas_programadas || 0}</div></div>`;
        html += `<div class="detail-item"><div class="label">Citas Completadas</div><div class="value">${r.citas_completadas || 0}</div></div>`;
        html += `<div class="detail-item"><div class="label">Comidas Registradas</div><div class="value">${r.comidas_registradas || 0}</div></div>`;
        html += `<div class="detail-item"><div class="label">Adherencia Medicinas</div><div class="value">${r.adherencia_medicinas || 0}%</div></div>`;
        html += `<div class="detail-item"><div class="label">Alertas Salud</div><div class="value">${r.alertas_salud || 0}</div></div>`;
        html += `<div class="detail-item"><div class="label">Estado General</div><div class="value">${badge(r.estado_general || 'estable')}</div></div>`;
        html += `</div>`;
        if (r.recomendaciones && r.recomendaciones.length > 0) {
          html += `<div style="margin-top:0.8rem"><div class="label" style="font-size:0.75rem;color:#6B7280;text-transform:uppercase;letter-spacing:0.5px;font-weight:600;margin-bottom:0.3rem">Recomendaciones</div><ul style="margin:0;padding-left:1.2rem;color:#374151">`;
          r.recomendaciones.forEach(rec => { html += `<li>${rec}</li>`; });
          html += `</ul></div>`;
        }
        html += `</div>`;
      });
    } else {
      html = `<div class="empty-state"><p>No hay reportes medicos</p></div>`;
    }
    body.innerHTML = html;
  } catch (err) {
    body.innerHTML = `<div class="empty-state"><p>Error: ${err.message}</p></div>`;
  }
}
