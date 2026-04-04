package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Prospector</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#4a7ec9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.stats-row{display:grid;grid-template-columns:repeat(3,1fr);gap:.6rem;padding:1rem 1.5rem;border-bottom:1px solid var(--bg3)}
.stat{text-align:center}.stat-val{font-family:var(--mono);font-size:1.3rem;color:var(--cream)}.stat-label{font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px}
.pipeline{display:grid;grid-template-columns:repeat(5,1fr);gap:.4rem;padding:1rem;overflow-x:auto}
@media(max-width:800px){.pipeline{grid-template-columns:repeat(5,200px)}}
.stage{background:var(--bg2);border:1px solid var(--bg3);min-height:300px}
.stage-hdr{font-family:var(--mono);font-size:.6rem;text-transform:uppercase;letter-spacing:1px;padding:.5rem .6rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between}
.stage-body{padding:.3rem}
.deal{background:var(--bg);border:1px solid var(--bg3);padding:.5rem .6rem;margin-bottom:.3rem;cursor:pointer;transition:border-color .15s}
.deal:hover{border-color:var(--leather)}
.deal-name{font-family:var(--mono);font-size:.72rem;margin-bottom:.1rem}
.deal-co{font-size:.68rem;color:var(--cd)}
.deal-val{font-family:var(--mono);font-size:.7rem;color:var(--gold);margin-top:.2rem}
.deal-prob{font-family:var(--mono);font-size:.55rem;color:var(--cm)}
.s-lead .stage-hdr{color:var(--cm)}.s-qualified .stage-hdr{color:var(--blue)}.s-proposal .stage-hdr{color:var(--leather)}
.s-negotiation .stage-hdr{color:var(--gold)}.s-won .stage-hdr{color:var(--green)}.s-lost .stage-hdr{color:var(--red)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.toolbar{padding:.6rem 1.5rem;border-bottom:1px solid var(--bg3)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem}
.fr textarea{min-height:40px;resize:vertical}
.actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:1.5rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>PROSPECTOR</h1><button class="btn btn-primary" onclick="openForm()" style="font-size:.6rem">+ New Deal</button></div>
<div class="stats-row" id="statsRow"></div>
<div class="pipeline" id="pipeline"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api',stages=['lead','qualified','proposal','negotiation','won'];
let deals=[];
async function load(){const[d,s]=await Promise.all([fetch(A+'/deals').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
deals=d.deals||[];
document.getElementById('statsRow').innerHTML='<div class="stat"><div class="stat-val">'+deals.length+'</div><div class="stat-label">Deals</div></div><div class="stat"><div class="stat-val">$'+fmt(s.pipeline_value||0)+'</div><div class="stat-label">Pipeline</div></div><div class="stat"><div class="stat-val">$'+fmt(s.won_value||0)+'</div><div class="stat-label">Won</div></div>';
render();}
function render(){const p=document.getElementById('pipeline');
let h='';stages.forEach(st=>{
const sd=(deals||[]).filter(d=>d.stage===st);
const val=sd.reduce((s,d)=>s+d.value,0);
h+='<div class="stage s-'+st+'"><div class="stage-hdr"><span>'+st.charAt(0).toUpperCase()+st.slice(1)+'</span><span>$'+fmt(val)+'</span></div><div class="stage-body">';
if(sd.length){sd.forEach(d=>{h+='<div class="deal" onclick="openEdit(\''+d.id+'\')"><div class="deal-name">'+esc(d.name)+'</div>';if(d.company)h+='<div class="deal-co">'+esc(d.company)+'</div>';h+='<div class="deal-val">$'+fmt(d.value)+'</div><div class="deal-prob">'+d.probability+'% · '+(d.close_date||'no date')+'</div></div>';});}
else{h+='<div class="empty">No deals</div>';}
h+='</div></div>';});
p.innerHTML=h;}
function openForm(){document.getElementById('mdl').innerHTML='<h2>New Deal</h2><div class="fr"><label>Deal Name</label><input id="f-name" placeholder="e.g. Enterprise license"></div><div class="fr"><label>Company</label><input id="f-co" placeholder="Acme Corp"></div><div class="fr"><label>Contact Name</label><input id="f-cn"></div><div class="fr"><label>Contact Email</label><input id="f-ce" type="email"></div><div class="fr"><label>Value ($)</label><input id="f-val" type="number" value="0"></div><div class="fr"><label>Stage</label><select id="f-stage"><option value="lead">Lead</option><option value="qualified">Qualified</option><option value="proposal">Proposal</option><option value="negotiation">Negotiation</option><option value="won">Won</option></select></div><div class="fr"><label>Probability (%)</label><input id="f-prob" type="number" min="0" max="100" value="10"></div><div class="fr"><label>Close Date</label><input id="f-date" type="date"></div><div class="fr"><label>Notes</label><textarea id="f-notes"></textarea></div><div class="actions"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-primary" onclick="submitDeal()">Create</button></div>';
document.getElementById('mbg').classList.add('open');}
async function submitDeal(){await fetch(A+'/deals',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-name').value,company:document.getElementById('f-co').value,contact_name:document.getElementById('f-cn').value,contact_email:document.getElementById('f-ce').value,value:parseInt(document.getElementById('f-val').value)||0,stage:document.getElementById('f-stage').value,probability:parseInt(document.getElementById('f-prob').value)||0,close_date:document.getElementById('f-date').value,notes:document.getElementById('f-notes').value})});cm();load();}
function openEdit(id){const d=(deals||[]).find(d=>d.id===id);if(!d)return;
document.getElementById('mdl').innerHTML='<h2>Edit Deal</h2><div class="fr"><label>Name</label><input id="e-name" value="'+esc(d.name)+'"></div><div class="fr"><label>Company</label><input id="e-co" value="'+esc(d.company||'')+'"></div><div class="fr"><label>Value ($)</label><input id="e-val" type="number" value="'+d.value+'"></div><div class="fr"><label>Stage</label><select id="e-stage">'+stages.concat(['lost']).map(s=>'<option value="'+s+'"'+(d.stage===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>').join('')+'</select></div><div class="fr"><label>Probability (%)</label><input id="e-prob" type="number" value="'+d.probability+'"></div><div class="fr"><label>Notes</label><textarea id="e-notes">'+(d.notes||'')+'</textarea></div><div class="actions"><button class="btn" style="color:var(--red)" onclick="delDeal(\''+id+'\')">Delete</button><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-primary" onclick="saveDeal(\''+id+'\')">Save</button></div>';
document.getElementById('mbg').classList.add('open');}
async function saveDeal(id){await fetch(A+'/deals/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('e-name').value,company:document.getElementById('e-co').value,value:parseInt(document.getElementById('e-val').value)||0,stage:document.getElementById('e-stage').value,probability:parseInt(document.getElementById('e-prob').value)||0,notes:document.getElementById('e-notes').value})});cm();load();}
async function delDeal(id){if(confirm('Delete?')){await fetch(A+'/deals/'+id,{method:'DELETE'});cm();load();}}
function cm(){document.getElementById('mbg').classList.remove('open');}
function fmt(n){if(n>=1000000)return(n/1000000).toFixed(1)+'M';if(n>=1000)return(n/1000).toFixed(0)+'k';return n.toString();}
function esc(s){if(!s)return'';return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');}
load();
</script></body></html>`
