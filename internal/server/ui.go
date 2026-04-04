package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Prospector</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.header{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.header h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.pipeline-stats{font-family:var(--mono);font-size:.7rem;color:var(--cm)}
.pipeline-stats .val{color:var(--gold)}
.content{padding:1rem;overflow-x:auto}
.pipeline{display:flex;gap:.6rem;min-width:900px}
.stage{flex:1;min-width:170px;background:var(--bg2);border:1px solid var(--bg3)}
.stage-header{padding:.6rem .8rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.stage-name{font-family:var(--mono);font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px}
.stage-total{font-family:var(--mono);font-size:.6rem;color:var(--gold)}
.deal{border-bottom:1px solid var(--bg3);padding:.5rem .8rem;cursor:pointer;transition:background .1s}
.deal:hover{background:var(--bg)}
.deal-name{font-family:var(--mono);font-size:.72rem;margin-bottom:.1rem}
.deal-company{font-size:.7rem;color:var(--cd)}
.deal-meta{display:flex;justify-content:space-between;margin-top:.2rem;font-family:var(--mono);font-size:.55rem;color:var(--cm)}
.deal-value{color:var(--gold)}
.deal-prob{color:var(--green)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-primary:hover{opacity:.85}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:450px;max-width:90vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.form-row{margin-bottom:.6rem}
.form-row label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.form-row input,.form-row select,.form-row textarea{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem}
.form-row textarea{min-height:50px;resize:vertical}
.actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
.empty-stage{padding:1rem;text-align:center;color:var(--cm);font-style:italic;font-size:.7rem}
</style></head><body>
<div class="header"><h1>PROSPECTOR</h1><div style="display:flex;gap:.8rem;align-items:center"><div class="pipeline-stats" id="stats"></div><button class="btn btn-primary" onclick="openForm()">+ New Deal</button></div></div>
<div class="content"><div class="pipeline" id="pipeline"></div></div>
<div class="modal-bg" id="modalBg" onclick="if(event.target===this)closeModal()"><div class="modal" id="modal"></div></div>

<script>
const API='/api';
const STAGES=['lead','qualified','proposal','negotiation','won','lost'];
const STAGE_LABELS={'lead':'Lead','qualified':'Qualified','proposal':'Proposal','negotiation':'Negotiation','won':'Won','lost':'Lost'};
let deals=[];

async function load(){
  const r=await fetch(API+'/deals').then(r=>r.json());
  deals=r.deals||[];render();
}

function render(){
  const totalValue=(deals||[]).reduce((s,d)=>s+(d.stage!=='lost'?d.value:0),0);
  const weighted=(deals||[]).reduce((s,d)=>s+(d.stage!=='lost'&&d.stage!=='won'?d.value*d.probability/100:d.stage==='won'?d.value:0),0);
  document.getElementById('stats').innerHTML='Pipeline: <span class="val">$'+fmt(totalValue)+'</span> &middot; Weighted: <span class="val">$'+fmt(weighted)+'</span>';

  let h='';
  STAGES.forEach(stage=>{
    const stageDials=(deals||[]).filter(d=>d.stage===stage);
    const stageTotal=stageDials.reduce((s,d)=>s+d.value,0);
    h+='<div class="stage"><div class="stage-header"><span class="stage-name">'+STAGE_LABELS[stage]+' ('+stageDials.length+')</span><span class="stage-total">$'+fmt(stageTotal)+'</span></div>';
    if(!stageDials.length)h+='<div class="empty-stage">No deals</div>';
    stageDials.forEach(d=>{
      h+='<div class="deal" onclick="openEdit(\''+d.id+'\')"><div class="deal-name">'+esc(d.name)+'</div>';
      if(d.company)h+='<div class="deal-company">'+esc(d.company)+'</div>';
      h+='<div class="deal-meta"><span class="deal-value">$'+fmt(d.value)+'</span><span class="deal-prob">'+d.probability+'%</span></div></div>';
    });
    h+='</div>';
  });
  document.getElementById('pipeline').innerHTML=h;
}

function openForm(deal){
  const d=deal||{name:'',company:'',contact_name:'',contact_email:'',value:0,stage:'lead',probability:10,close_date:'',notes:''};
  const isEdit=!!deal;
  document.getElementById('modal').innerHTML='<h2>'+(isEdit?'Edit':'New')+' Deal</h2><div class="form-row"><label>Deal name</label><input id="f-name" value="'+esc(d.name)+'"></div><div class="form-row"><label>Company</label><input id="f-company" value="'+esc(d.company)+'"></div><div class="form-row"><label>Contact name</label><input id="f-contact" value="'+esc(d.contact_name)+'"></div><div class="form-row"><label>Contact email</label><input id="f-email" value="'+esc(d.contact_email)+'"></div><div class="form-row"><label>Value ($)</label><input id="f-value" type="number" value="'+d.value+'"></div><div class="form-row"><label>Stage</label><select id="f-stage">'+STAGES.map(s=>'<option value="'+s+'"'+(s===d.stage?' selected':'')+'>'+STAGE_LABELS[s]+'</option>').join('')+'</select></div><div class="form-row"><label>Probability (%)</label><input id="f-prob" type="number" min="0" max="100" value="'+d.probability+'"></div><div class="form-row"><label>Expected close</label><input id="f-close" type="date" value="'+esc(d.close_date)+'"></div><div class="form-row"><label>Notes</label><textarea id="f-notes">'+esc(d.notes||'')+'</textarea></div><div class="actions">'+(isEdit?'<button class="btn" onclick="delDeal(\''+d.id+'\')" style="color:var(--red)">Delete</button>':'')+'<button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-primary" onclick="'+(isEdit?'updateDeal(\''+d.id+'\')':'createDeal()')+'">Save</button></div>';
  document.getElementById('modalBg').classList.add('open');
}

function openEdit(id){const d=(deals||[]).find(x=>x.id===id);if(d)openForm(d);}

async function createDeal(){
  await fetch(API+'/deals',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(formData())});closeModal();load();
}
async function updateDeal(id){
  await fetch(API+'/deals/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(formData())});closeModal();load();
}
async function delDeal(id){if(confirm('Delete?')){await fetch(API+'/deals/'+id,{method:'DELETE'});closeModal();load();}}

function formData(){return{name:document.getElementById('f-name').value,company:document.getElementById('f-company').value,contact_name:document.getElementById('f-contact').value,contact_email:document.getElementById('f-email').value,value:parseInt(document.getElementById('f-value').value)||0,stage:document.getElementById('f-stage').value,probability:parseInt(document.getElementById('f-prob').value)||0,close_date:document.getElementById('f-close').value,notes:document.getElementById('f-notes').value};}

function closeModal(){document.getElementById('modalBg').classList.remove('open');}
function fmt(n){return n>=1000?(n/1000).toFixed(n>=10000?0:1)+'k':n.toString();}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
