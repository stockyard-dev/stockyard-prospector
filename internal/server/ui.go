package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Prospector</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.hdr-stats{font-family:var(--mono);font-size:.7rem;color:var(--cm)}
.hdr-stats strong{color:var(--gold)}
.pipeline{display:grid;grid-template-columns:repeat(5,1fr);gap:.5rem;padding:1rem;min-height:calc(100vh - 55px);overflow-x:auto}
@media(max-width:900px){.pipeline{grid-template-columns:repeat(5,200px)}}
.stage{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem}
.stage-hdr{font-family:var(--mono);font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px;margin-bottom:.6rem;display:flex;justify-content:space-between}
.stage-count{background:var(--bg3);padding:.1rem .4rem;font-size:.55rem;color:var(--cm)}
.deal{background:var(--bg);border:1px solid var(--bg3);padding:.6rem;margin-bottom:.4rem;cursor:pointer;transition:border-color .15s}
.deal:hover{border-color:var(--leather)}
.deal-name{font-family:var(--mono);font-size:.75rem}
.deal-company{font-size:.7rem;color:var(--cm)}
.deal-value{font-family:var(--mono);font-size:.7rem;color:var(--gold);margin-top:.2rem}
.deal-prob{font-family:var(--mono);font-size:.55rem;color:var(--cm)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-p:hover{opacity:.85}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem}
.fr textarea{min-height:60px;resize:vertical}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.stage-btns{display:flex;gap:.3rem;margin-top:.5rem;flex-wrap:wrap}
.stage-btns .btn{font-size:.55rem;padding:.2rem .4rem}
</style></head><body>
<div class="hdr"><div><h1>PROSPECTOR</h1></div><div style="display:flex;gap:1rem;align-items:center"><div class="hdr-stats" id="st"></div><button class="btn btn-p" onclick="openDeal()">+ New Deal</button></div></div>
<div class="pipeline" id="board"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api',stages=['lead','qualified','proposal','negotiation','closed_won'];
const stageLabels={lead:'Lead',qualified:'Qualified',proposal:'Proposal',negotiation:'Negotiation',closed_won:'Closed Won'};
let deals=[];
async function ld(){const[d,s]=await Promise.all([fetch(A+'/deals').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);deals=d.deals||[];
document.getElementById('st').innerHTML='<strong>$'+fmt(s.total_value||0)+'</strong> pipeline &middot; <strong>$'+fmt(s.weighted_value||0)+'</strong> weighted &middot; '+s.total+' deals';rn();}
function rn(){
  let h='';
  stages.forEach(s=>{
    const items=deals.filter(d=>d.stage===s);
    const val=items.reduce((a,d)=>a+d.value,0);
    h+='<div class="stage"><div class="stage-hdr"><span>'+stageLabels[s]+'</span><span class="stage-count">'+items.length+' &middot; $'+fmt(val)+'</span></div>';
    items.forEach(d=>{
      h+='<div class="deal" onclick="openDetail(\''+d.id+'\')"><div class="deal-name">'+esc(d.name)+'</div>';
      if(d.company)h+='<div class="deal-company">'+esc(d.company)+'</div>';
      h+='<div class="deal-value">$'+fmt(d.value)+'<span class="deal-prob"> &middot; '+d.probability+'%</span></div></div>';
    });
    h+='</div>';
  });
  document.getElementById('board').innerHTML=h;
}
function openDeal(){
  document.getElementById('mdl').innerHTML='<h2>New Deal</h2><div class="fr"><label>Deal Name</label><input id="dn" placeholder="e.g. Acme Corp Enterprise"></div><div class="fr"><label>Company</label><input id="dc" placeholder="e.g. Acme Corp"></div><div class="fr"><label>Contact Name</label><input id="dcn"></div><div class="fr"><label>Contact Email</label><input id="dce" type="email"></div><div class="fr"><label>Value ($)</label><input id="dv" type="number" value="0"></div><div class="fr"><label>Stage</label><select id="ds">'+stages.map(s=>'<option value="'+s+'">'+stageLabels[s]+'</option>').join('')+'</select></div><div class="fr"><label>Probability (%)</label><input id="dp" type="number" value="10" min="0" max="100"></div><div class="fr"><label>Close Date</label><input id="dd" type="date"></div><div class="fr"><label>Notes</label><textarea id="dno"></textarea></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sDeal()">Create</button></div>';
  document.getElementById('mbg').classList.add('open');
}
async function sDeal(){await fetch(A+'/deals',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('dn').value,company:document.getElementById('dc').value,contact_name:document.getElementById('dcn').value,contact_email:document.getElementById('dce').value,value:parseInt(document.getElementById('dv').value)||0,stage:document.getElementById('ds').value,probability:parseInt(document.getElementById('dp').value)||0,close_date:document.getElementById('dd').value,notes:document.getElementById('dno').value})});cm();ld();}
function openDetail(id){
  const d=deals.find(x=>x.id===id);if(!d)return;
  let h='<h2>'+esc(d.name)+'</h2>';
  if(d.company)h+='<p style="color:var(--cd);margin-bottom:.3rem">'+esc(d.company)+'</p>';
  if(d.contact_name)h+='<p style="font-family:var(--mono);font-size:.72rem;color:var(--cm)">'+esc(d.contact_name)+(d.contact_email?' &lt;'+esc(d.contact_email)+'&gt;':'')+'</p>';
  h+='<div style="font-family:var(--mono);font-size:.8rem;color:var(--gold);margin:.5rem 0">$'+fmt(d.value)+' &middot; '+d.probability+'% probability</div>';
  if(d.close_date)h+='<div style="font-family:var(--mono);font-size:.7rem;color:var(--cm)">Close: '+d.close_date+'</div>';
  if(d.notes)h+='<div style="font-size:.82rem;color:var(--cd);margin:.5rem 0;padding:.5rem;background:var(--bg);border:1px solid var(--bg3)">'+esc(d.notes)+'</div>';
  h+='<div style="font-family:var(--mono);font-size:.6rem;color:var(--cm);margin:.5rem 0">Move to:</div><div class="stage-btns">';
  stages.forEach(s=>{h+='<button class="btn'+(d.stage===s?' btn-p':'')+'" onclick="ms(\''+id+'\',\''+s+'\')">'+stageLabels[s]+'</button>';});
  h+='</div><div class="acts" style="margin-top:1rem"><button class="btn" style="color:var(--red)" onclick="dd(\''+id+'\')">Delete</button><button class="btn" onclick="cm()">Close</button></div>';
  document.getElementById('mdl').innerHTML=h;document.getElementById('mbg').classList.add('open');
}
async function ms(id,stage){await fetch(A+'/deals/'+id+'/stage',{method:'PATCH',headers:{'Content-Type':'application/json'},body:JSON.stringify({stage})});cm();ld();}
async function dd(id){if(confirm('Delete deal?')){await fetch(A+'/deals/'+id,{method:'DELETE'});cm();ld();}}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
function fmt(n){return n>=1000000?(n/1000000).toFixed(1)+'M':n>=1000?(n/1000).toFixed(0)+'k':n.toString();}
ld();
</script></body></html>`
