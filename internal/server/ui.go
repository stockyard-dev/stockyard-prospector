package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Prospector</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:1280px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-v.green{color:var(--green)}
.st-v.blue{color:var(--blue)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.pipeline{display:grid;grid-template-columns:repeat(5,minmax(0,1fr));gap:.5rem;overflow-x:auto}
.col{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;min-height:300px;min-width:0}
.col-hdr{font-size:.55rem;text-transform:uppercase;letter-spacing:1px;margin-bottom:.6rem;display:flex;justify-content:space-between;color:var(--cm);padding-bottom:.4rem;border-bottom:1px solid var(--bg3)}
.col-val{color:var(--gold);font-weight:700}
.deal{background:var(--bg);border:1px solid var(--bg3);padding:.5rem .6rem;margin-bottom:.4rem;font-size:.65rem;transition:border-color .15s;border-left:3px solid transparent}
.deal:hover{border-color:var(--leather)}
.deal-name{color:var(--cream);font-size:.72rem;font-weight:700;margin-bottom:.15rem}
.deal-co{color:var(--cd);font-size:.6rem;margin-bottom:.2rem}
.deal-val{color:var(--gold);font-size:.7rem;font-weight:700;margin-top:.2rem}
.deal-val .prob{color:var(--cm);font-weight:400;font-size:.55rem;margin-left:.3rem}
.deal-meta{font-size:.52rem;color:var(--cm);margin-top:.3rem;display:flex;flex-direction:column;gap:.1rem}
.deal-extra{font-size:.5rem;color:var(--cd);margin-top:.3rem;padding-top:.25rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.1rem}
.deal-extra-row{display:flex;gap:.3rem}
.deal-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.3px;min-width:70px}
.deal-extra-val{color:var(--cream)}
.deal-acts{display:flex;gap:.25rem;margin-top:.4rem;flex-wrap:wrap}
.btn{font-size:.55rem;padding:.18rem .35rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .15s;font-family:var(--mono)}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff;font-size:.65rem;padding:.3rem .6rem}
.btn-fwd{color:var(--green);border-color:var(--bg3)}
.btn-fwd:hover{border-color:var(--green)}
.btn-lost{color:var(--red)}
.btn-lost:hover{border-color:var(--red)}
.empty-col{text-align:center;padding:1.5rem .5rem;color:var(--cm);font-style:italic;font-size:.6rem}
.lost-summary{padding:.6rem 1rem;font-size:.6rem;color:var(--cm);text-align:center;margin-top:.6rem;border:1px dashed var(--bg3);background:var(--bg2)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
@media(max-width:1000px){.pipeline{grid-template-columns:repeat(3,minmax(0,1fr))}.stats{grid-template-columns:repeat(2,1fr)}}
@media(max-width:600px){.pipeline{grid-template-columns:1fr}.stats{grid-template-columns:1fr 1fr}.row2{grid-template-columns:1fr}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> PROSPECTOR</h1>
<button class="btn-p btn" onclick="openForm()">+ New Deal</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search name, company, contact..." oninput="render()">
</div>
<div class="pipeline" id="pipeline"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='deals';

// The five active stages plus 'lost'. STAGES is the kanban column order
// (lost is shown as a summary line below the columns, not as its own column).
var STAGES=['lead','qualified','proposal','negotiation','won'];

// Field defs drive the form, the cards, and the submit body.
// value is integer dollars (not cents) — for backward compat with the
// original schema. probability is 0-100.
var fields=[
{name:'name',label:'Deal Name',type:'text',required:true,placeholder:'Acme Corp - Enterprise'},
{name:'company',label:'Company',type:'text'},
{name:'contact_name',label:'Contact Name',type:'text'},
{name:'contact_email',label:'Contact Email',type:'email'},
{name:'value',label:'Value ($)',type:'integer',placeholder:'10000'},
{name:'stage',label:'Stage',type:'select',required:true,options:['lead','qualified','proposal','negotiation','won','lost']},
{name:'probability',label:'Probability (%)',type:'integer',placeholder:'50'},
{name:'close_date',label:'Expected Close',type:'date'},
{name:'notes',label:'Notes',type:'textarea'}
];

var deals=[],editId=null;

// ─── Money helpers ────────────────────────────────────────────────
// Value is integer dollars. fmtMoney compacts large amounts: $1.5M, $250k.

function fmtMoney(dollars){
var n=parseInt(dollars||0,10);
if(isNaN(n))return'$0';
var neg=n<0;
n=Math.abs(n);
var s;
if(n>=1000000)s='$'+(n/1000000).toFixed(n>=10000000?0:1)+'M';
else if(n>=10000)s='$'+Math.round(n/1000)+'k';
else if(n>=1000)s='$'+(n/1000).toFixed(1)+'k';
else s='$'+n.toLocaleString();
return neg?'-'+s:s;
}

function fmtMoneyFull(dollars){
var n=parseInt(dollars||0,10);
if(isNaN(n))return'$0';
var neg=n<0;
n=Math.abs(n);
var s='$'+n.toLocaleString();
return neg?'-'+s:s;
}

function fmtDate(s){
if(!s)return'';
try{return new Date(s).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}catch(e){return s}
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resp=await fetch(A+'/'+RESOURCE).then(function(r){return r.json()});
var list=resp[RESOURCE]||[];
try{
var extras=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
list.forEach(function(d){
var ex=extras[d.id];
if(!ex)return;
Object.keys(ex).forEach(function(k){if(d[k]===undefined)d[k]=ex[k]});
});
}catch(e){}
deals=list;
}catch(e){
console.error('load failed',e);
deals=[];
}
renderStats();
render();
}

function renderStats(){
var total=deals.length;
var pipelineVal=0;
var wonVal=0;
var weightedVal=0;
deals.forEach(function(d){
var v=parseInt(d.value||0,10);
var p=parseInt(d.probability||0,10);
if(d.stage==='won')wonVal+=v;
else if(d.stage!=='lost'){pipelineVal+=v;weightedVal+=v*p/100}
});
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Total Deals</div></div>'+
'<div class="st"><div class="st-v">'+fmtMoney(pipelineVal)+'</div><div class="st-l">Pipeline</div></div>'+
'<div class="st"><div class="st-v blue">'+fmtMoney(Math.round(weightedVal))+'</div><div class="st-l">Weighted</div></div>'+
'<div class="st"><div class="st-v green">'+fmtMoney(wonVal)+'</div><div class="st-l">Won</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var filtered=deals;
if(q)filtered=deals.filter(function(d){
return(d.name||'').toLowerCase().includes(q)||
       (d.company||'').toLowerCase().includes(q)||
       (d.contact_name||'').toLowerCase().includes(q)||
       (d.contact_email||'').toLowerCase().includes(q);
});

var h='';
STAGES.forEach(function(stage){
var sd=filtered.filter(function(d){return d.stage===stage});
var val=sd.reduce(function(s,d){return s+parseInt(d.value||0,10)},0);
h+='<div class="col">';
h+='<div class="col-hdr"><span>'+stage.toUpperCase()+' ('+sd.length+')</span><span class="col-val">'+fmtMoney(val)+'</span></div>';
if(!sd.length){
h+='<div class="empty-col">No deals</div>';
}else{
sd.forEach(function(d){h+=dealHTML(d)});
}
h+='</div>';
});
document.getElementById('pipeline').innerHTML=h;

var lost=filtered.filter(function(d){return d.stage==='lost'});
if(lost.length){
var lostVal=lost.reduce(function(s,d){return s+parseInt(d.value||0,10)},0);
var lostHTML='<div class="lost-summary">'+lost.length+' lost deal'+(lost.length!==1?'s':'')+' &middot; '+fmtMoney(lostVal)+' total &middot; <a href="#" onclick="event.preventDefault();showLostList()" style="color:var(--leather)">view</a></div>';
document.getElementById('pipeline').insertAdjacentHTML('afterend',lostHTML);
}
}

function dealHTML(d){
var si=STAGES.indexOf(d.stage);
var next=si>=0&&si<STAGES.length-1?STAGES[si+1]:null;

var h='<div class="deal">';
h+='<div class="deal-name">'+esc(d.name)+'</div>';
if(d.company)h+='<div class="deal-co">'+esc(d.company)+'</div>';
h+='<div class="deal-val">'+fmtMoneyFull(d.value);
if(d.probability)h+='<span class="prob">'+esc(String(d.probability))+'%</span>';
h+='</div>';

h+='<div class="deal-meta">';
if(d.contact_name)h+='<span>'+esc(d.contact_name)+'</span>';
if(d.close_date)h+='<span>close: '+esc(fmtDate(d.close_date))+'</span>';
h+='</div>';

// Custom fields from personalization
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var v=d[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="deal-extra-row">';
customRows+='<span class="deal-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="deal-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="deal-extra">'+customRows+'</div>';

h+='<div class="deal-acts">';
if(next)h+='<button class="btn btn-fwd" onclick="moveStage(\''+d.id+'\',\''+next+'\')">→ '+next+'</button>';
if(d.stage!=='won'&&d.stage!=='lost')h+='<button class="btn btn-lost" onclick="moveStage(\''+d.id+'\',\'lost\')">Lost</button>';
h+='<button class="btn" onclick="openEdit(\''+d.id+'\')">Edit</button>';
h+='<button class="btn" onclick="del(\''+d.id+'\')" style="color:var(--cm)">&#10005;</button>';
h+='</div>';

h+='</div>';
return h;
}

function showLostList(){
var lost=deals.filter(function(d){return d.stage==='lost'});
if(!lost.length)return;
var html='<h2>LOST DEALS</h2>';
lost.forEach(function(d){
html+='<div style="padding:.5rem;border:1px solid var(--bg3);margin-bottom:.4rem">';
html+='<div style="font-weight:700;font-size:.75rem">'+esc(d.name)+'</div>';
if(d.company)html+='<div style="font-size:.65rem;color:var(--cd)">'+esc(d.company)+'</div>';
html+='<div style="font-size:.65rem;color:var(--gold);margin-top:.2rem">'+fmtMoneyFull(d.value)+'</div>';
html+='<div style="display:flex;gap:.3rem;margin-top:.4rem">';
html+='<button class="btn" onclick="moveStage(\''+d.id+'\',\'lead\');closeModal()">Re-open as Lead</button>';
html+='<button class="btn" onclick="openEdit(\''+d.id+'\')">Edit</button>';
html+='<button class="btn" onclick="if(confirm(\\'Permanently delete?\\'))del(\''+d.id+'\')" style="color:var(--red)">Delete</button>';
html+='</div></div>';
});
html+='<div class="acts"><button class="btn" onclick="closeModal()">Close</button></div>';
document.getElementById('mdl').innerHTML=html;
document.getElementById('mbg').classList.add('open');
}

// ─── Form ─────────────────────────────────────────────────────────

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph='';
if(f.placeholder)ph=' placeholder="'+esc(f.placeholder)+'"';
else if(f.name==='name'&&window._placeholderName)ph=' placeholder="'+esc(window._placeholderName)+'"';

var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=(typeof o==='string')?(o.charAt(0).toUpperCase()+o.slice(1)):String(o);
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="3"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='checkbox'){
h+='<input type="checkbox" id="f-'+f.name+'"'+(v?' checked':'')+' style="width:auto">';
}else if(f.type==='integer'||f.type==='number'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}

h+='</div>';
return h;
}

function formHTML(deal){
var i=deal||{};
var isEdit=!!deal;
var h='<h2>'+(isEdit?'EDIT DEAL':'NEW DEAL')+'</h2>';

// Name on its own (most important)
h+=fieldHTML(fieldByName('name'),i.name);

// Company on its own
h+=fieldHTML(fieldByName('company'),i.company);

// Contact name + email pair
h+='<div class="row2">'+fieldHTML(fieldByName('contact_name'),i.contact_name)+fieldHTML(fieldByName('contact_email'),i.contact_email)+'</div>';

// Value + probability pair
h+='<div class="row2">'+fieldHTML(fieldByName('value'),i.value)+fieldHTML(fieldByName('probability'),i.probability)+'</div>';

// Stage + close date pair
h+='<div class="row2">'+fieldHTML(fieldByName('stage'),i.stage||'lead')+fieldHTML(fieldByName('close_date'),i.close_date)+'</div>';

// Notes on its own
h+=fieldHTML(fieldByName('notes'),i.notes);

// Custom fields injected by personalization
var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var sectionLabel=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(sectionLabel)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,i[f.name])});
h+='</div>';
}

h+='<div class="acts">';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn-p btn" onclick="submit()">'+(isEdit?'Save':'Create Deal')+'</button>';
h+='</div>';
return h;
}

function openForm(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var x=null;
for(var j=0;j<deals.length;j++){if(deals[j].id===id){x=deals[j];break}}
if(!x)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(x);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

// ─── Submit ───────────────────────────────────────────────────────

async function submit(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Deal name is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='checkbox')val=el.checked?1:0;
else if(f.type==='integer')val=parseInt(el.value,10)||0;
else if(f.type==='number')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/'+RESOURCE+'/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/'+RESOURCE,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Save failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}

closeModal();
load();
}

// Quick stage move via PATCH endpoint (kanban arrow buttons)
async function moveStage(id,stage){
try{
await fetch(A+'/'+RESOURCE+'/'+id+'/stage',{method:'PATCH',headers:{'Content-Type':'application/json'},body:JSON.stringify({stage:stage})});
load();
}catch(e){alert('Move failed: '+e.message)}
}

async function del(id){
if(!confirm('Delete this deal?'))return;
await fetch(A+'/'+RESOURCE+'/'+id,{method:'DELETE'});
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.placeholder_name)window._placeholderName=cfg.placeholder_name;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
