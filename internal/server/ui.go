package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Prospector</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.6rem;padding:1rem 1.5rem;max-width:1200px;margin:0 auto}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem}.st-l{font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.1rem}
.pipeline{display:grid;grid-template-columns:repeat(5,1fr);gap:.5rem;padding:0 1rem 1rem;max-width:1200px;margin:0 auto;overflow-x:auto}
@media(max-width:800px){.pipeline{grid-template-columns:1fr}}
.col{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;min-height:200px}
.col-hdr{font-size:.6rem;text-transform:uppercase;letter-spacing:1px;margin-bottom:.6rem;display:flex;justify-content:space-between;color:var(--cm)}
.col-val{color:var(--gold)}
.deal{background:var(--bg);border:1px solid var(--bg3);padding:.5rem .7rem;margin-bottom:.4rem;cursor:pointer;transition:border-color .15s;font-size:.7rem}
.deal:hover{border-color:var(--leather)}
.deal-name{color:var(--cream);font-size:.75rem}
.deal-co{color:var(--cd);font-size:.65rem}
.deal-val{color:var(--gold);font-size:.7rem;margin-top:.2rem}
.deal-meta{font-size:.55rem;color:var(--cm);margin-top:.2rem;display:flex;gap:.5rem}
.deal-acts{display:flex;gap:.3rem;margin-top:.3rem}
.btn{font-size:.55rem;padding:.2rem .4rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg);font-size:.6rem;padding:.25rem .6rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
</style></head><body>
<div class="hdr"><h1>PROSPECTOR</h1><button class="btn btn-p" onclick="openForm()">+ New Deal</button></div>
<div class="stats" id="stats"></div>
<div class="pipeline" id="pipeline"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api',STAGES=['lead','qualified','proposal','negotiation','won'];
let deals=[];
async function load(){
  const[d,s]=await Promise.all([fetch(A+'/deals').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
  deals=d.deals||[];
  document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+(s.total||0)+'</div><div class="st-l">Deals</div></div><div class="st"><div class="st-v">$'+fmt(s.pipeline_value||0)+'</div><div class="st-l">Pipeline</div></div><div class="st"><div class="st-v" style="color:var(--green)">$'+fmt(s.won_value||0)+'</div><div class="st-l">Won</div></div>';
  render();
}
function render(){
  let h='';
  STAGES.forEach(stage=>{
    const sd=deals.filter(d=>d.stage===stage);
    const val=sd.reduce((s,d)=>s+d.value,0);
    h+='<div class="col"><div class="col-hdr"><span>'+stage.toUpperCase()+' ('+sd.length+')</span><span class="col-val">$'+fmt(val)+'</span></div>';
    sd.forEach(d=>{
      const si=STAGES.indexOf(d.stage);
      const next=si<STAGES.length-1?STAGES[si+1]:null;
      h+='<div class="deal"><div class="deal-name">'+esc(d.name)+'</div>';
      if(d.company)h+='<div class="deal-co">'+esc(d.company)+'</div>';
      h+='<div class="deal-val">$'+fmt(d.value)+(d.probability?' · '+d.probability+'%':'')+'</div>';
      h+='<div class="deal-meta">';
      if(d.contact_name)h+='<span>'+esc(d.contact_name)+'</span>';
      if(d.close_date)h+='<span>close: '+d.close_date+'</span>';
      h+='</div>';
      h+='<div class="deal-acts">';
      if(next)h+='<button class="btn" onclick="mv(\''+d.id+'\',\''+next+'\')">→ '+next+'</button>';
      if(d.stage!=='won')h+='<button class="btn" onclick="mv(\''+d.id+'\',\'lost\')" style="color:var(--red)">Lost</button>';
      h+='<button class="btn" onclick="del(\''+d.id+'\')" style="color:var(--cm)">✕</button>';
      h+='</div></div>';
    });
    h+='</div>';
  });
  // Show lost deals count
  const lost=deals.filter(d=>d.stage==='lost');
  if(lost.length)h+='<div style="padding:.5rem 0;font-size:.6rem;color:var(--cm);grid-column:1/-1">'+lost.length+' lost deals hidden</div>';
  document.getElementById('pipeline').innerHTML=h;
}
async function mv(id,stage){await fetch(A+'/deals/'+id+'/stage',{method:'PATCH',headers:{'Content-Type':'application/json'},body:JSON.stringify({stage})});load();}
async function del(id){if(confirm('Delete?')){await fetch(A+'/deals/'+id,{method:'DELETE'});load();}}
function openForm(){
  document.getElementById('mdl').innerHTML='<h2>New Deal</h2><div class="fr"><label>Deal Name</label><input id="f-n" placeholder="e.g. Acme Corp - Enterprise"></div><div class="fr"><label>Company</label><input id="f-co"></div><div class="fr"><label>Contact Name</label><input id="f-cn"></div><div class="fr"><label>Contact Email</label><input id="f-ce" type="email"></div><div class="fr"><label>Value ($)</label><input id="f-v" type="number" value="0"></div><div class="fr"><label>Stage</label><select id="f-s"><option value="lead">Lead</option><option value="qualified">Qualified</option><option value="proposal">Proposal</option><option value="negotiation">Negotiation</option></select></div><div class="fr"><label>Probability (%)</label><input id="f-p" type="number" min="0" max="100" value="10"></div><div class="fr"><label>Expected Close</label><input id="f-d" type="date"></div><div class="fr"><label>Notes</label><textarea id="f-nt" rows="2"></textarea></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Create</button></div>';
  document.getElementById('mbg').classList.add('open');
}
async function sub(){await fetch(A+'/deals',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-n').value,company:document.getElementById('f-co').value,contact_name:document.getElementById('f-cn').value,contact_email:document.getElementById('f-ce').value,value:parseInt(document.getElementById('f-v').value)||0,stage:document.getElementById('f-s').value,probability:parseInt(document.getElementById('f-p').value)||0,close_date:document.getElementById('f-d').value,notes:document.getElementById('f-nt').value})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function fmt(n){return n>=1000000?(n/1000000).toFixed(1)+'M':n>=1000?(n/1000).toFixed(0)+'k':n.toString();}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
