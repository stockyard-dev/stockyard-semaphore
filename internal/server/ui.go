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
<title>Semaphore</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--purple:#9d6bb8;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1100px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(5,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center;cursor:pointer;transition:.15s}
.st:hover{border-color:var(--leather)}
.st.active{border-color:var(--rust)}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-v.green{color:var(--green)}
.st-v.orange{color:var(--orange)}
.st-v.red{color:var(--red)}
.st-v.blue{color:var(--blue)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:.6rem}
.card{background:var(--bg2);border:1px solid var(--bg3);padding:.9rem 1rem;display:flex;flex-direction:column;gap:.4rem;transition:border-color .15s}
.card:hover{border-color:var(--leather)}
.card-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem;cursor:pointer}
.card-name{font-size:.85rem;font-weight:700;color:var(--cream)}
.card-role{font-size:.62rem;color:var(--cd);margin-top:.1rem}
.dot{width:10px;height:10px;border-radius:50%;flex-shrink:0;margin-top:5px}
.dot.available{background:var(--green)}
.dot.busy{background:var(--red)}
.dot.away{background:var(--orange)}
.dot.in_meeting{background:var(--purple)}
.dot.off{background:var(--cm)}
.card-status{font-size:.65rem;color:var(--cd);font-style:italic;display:flex;align-items:center;gap:.4rem}
.card-meta{font-size:.55rem;color:var(--cm);display:flex;gap:.5rem;flex-wrap:wrap}
.card-meta a{color:var(--cd);text-decoration:none}
.card-meta a:hover{color:var(--rust)}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.available{border-color:var(--green);color:var(--green)}
.badge.busy{border-color:var(--red);color:var(--red)}
.badge.away{border-color:var(--orange);color:var(--orange)}
.badge.in_meeting{border-color:var(--purple);color:var(--purple)}
.badge.off{border-color:var(--cm);color:var(--cm)}
.quick-row{display:flex;gap:.2rem;margin-top:.3rem;flex-wrap:wrap}
.quick-btn{font-family:var(--mono);font-size:.55rem;padding:.2rem .35rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.quick-btn:hover{border-color:var(--leather);color:var(--cream)}
.card-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.card-extra-row{display:flex;gap:.4rem}
.card-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.card-extra-val{color:var(--cream)}

.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-del{color:var(--red);border-color:#3a1a1a}
.btn-del:hover{border-color:var(--red);color:var(--red)}

.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(3,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> SEMAPHORE</h1>
<button class="btn btn-p" onclick="openNew()">+ Add Member</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search name, email, role..." oninput="debouncedRender()">
<select class="filter-sel" id="role-filter" onchange="render()">
<option value="">All Roles</option>
</select>
</div>
<div id="grid" class="grid"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='team_members';

var fields=[
{name:'name',label:'Name',type:'text',required:true},
{name:'email',label:'Email',type:'email'},
{name:'role',label:'Role',type:'text'},
{name:'availability',label:'Availability',type:'select',options:['available','busy','away','in_meeting','off']},
{name:'status_message',label:'Status Message',type:'text'},
{name:'timezone',label:'Timezone',type:'text',placeholder:'e.g. America/Los_Angeles'},
{name:'until_time',label:'Until (auto-clear)',type:'datetime-local'}
];

var members=[],memberExtras={},editId=null,searchTimer=null,activeStatus='';

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
return d.toLocaleString('en-US',{month:'short',day:'numeric',hour:'numeric',minute:'2-digit'});
}catch(e){return s}
}

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/members').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
members=resps[0].members||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
memberExtras=ex||{};
members.forEach(function(m){
var x=memberExtras[m.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(m[k]===undefined)m[k]=x[k]});
});
}catch(e){memberExtras={}}

populateRoleFilter();
}catch(e){
console.error('load failed',e);
members=[];
}
render();
}

function populateRoleFilter(){
var sel=document.getElementById('role-filter');
if(!sel)return;
var current=sel.value;
var seen={};
var roles=[];
members.forEach(function(m){if(m.role&&!seen[m.role]){seen[m.role]=true;roles.push(m.role)}});
roles.sort();
sel.innerHTML='<option value="">All Roles</option>'+roles.map(function(r){return'<option value="'+esc(r)+'"'+(r===current?' selected':'')+'>'+esc(r)+'</option>'}).join('');
}

function renderStats(s){
var total=s.total||0;
var byA=s.by_availability||{};
document.getElementById('stats').innerHTML=
'<div class="st'+(activeStatus===''?' active':'')+'" onclick="filterStatus(\'\')"><div class="st-v">'+total+'</div><div class="st-l">Total</div></div>'+
'<div class="st'+(activeStatus==='available'?' active':'')+'" onclick="filterStatus(\'available\')"><div class="st-v green">'+(byA.available||0)+'</div><div class="st-l">Available</div></div>'+
'<div class="st'+(activeStatus==='busy'?' active':'')+'" onclick="filterStatus(\'busy\')"><div class="st-v red">'+(byA.busy||0)+'</div><div class="st-l">Busy</div></div>'+
'<div class="st'+(activeStatus==='in_meeting'?' active':'')+'" onclick="filterStatus(\'in_meeting\')"><div class="st-v blue">'+(byA.in_meeting||0)+'</div><div class="st-l">In Meeting</div></div>'+
'<div class="st'+(activeStatus==='away'||activeStatus==='off'?' active':'')+'" onclick="filterStatus(\'away\')"><div class="st-v orange">'+((byA.away||0)+(byA.off||0))+'</div><div class="st-l">Away/Off</div></div>';
}

function filterStatus(s){
activeStatus=s;
load();
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var rf=document.getElementById('role-filter').value;

var f=members;
if(q)f=f.filter(function(m){
return(m.name||'').toLowerCase().includes(q)||
(m.email||'').toLowerCase().includes(q)||
(m.role||'').toLowerCase().includes(q);
});
if(rf)f=f.filter(function(m){return m.role===rf});
if(activeStatus)f=f.filter(function(m){return m.availability===activeStatus});

if(!f.length){
var msg=window._emptyMsg||'No members match your filters.';
document.getElementById('grid').innerHTML='<div class="empty" style="grid-column:1/-1">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(m){h+=cardHTML(m)});
document.getElementById('grid').innerHTML=h;
}

function cardHTML(m){
var av=m.availability||'available';
var h='<div class="card">';
h+='<div class="card-top" onclick="openEdit(\''+esc(m.id)+'\')">';
h+='<div style="flex:1;min-width:0">';
h+='<div class="card-name">'+esc(m.name)+'</div>';
if(m.role)h+='<div class="card-role">'+esc(m.role)+'</div>';
h+='</div>';
h+='<div class="dot '+esc(av)+'" title="'+esc(av)+'"></div>';
h+='</div>';

h+='<div class="card-status">';
h+='<span class="badge '+esc(av)+'">'+esc(av.replace(/_/g,' '))+'</span>';
if(m.status_message)h+='<span>'+esc(m.status_message)+'</span>';
h+='</div>';

if(m.until_time)h+='<div class="card-meta" style="color:var(--orange)">until '+esc(fmtDate(m.until_time))+'</div>';

if(m.email||m.timezone){
h+='<div class="card-meta">';
if(m.email)h+='<a href="mailto:'+esc(m.email)+'" onclick="event.stopPropagation()">'+esc(m.email)+'</a>';
if(m.timezone)h+='<span>'+esc(m.timezone)+'</span>';
h+='</div>';
}

// Quick status buttons
h+='<div class="quick-row">';
['available','busy','in_meeting','away','off'].forEach(function(s){
if(s===av)return;
h+='<button class="quick-btn" onclick="quickSet(\''+esc(m.id)+'\',\''+s+'\')">'+s.replace(/_/g,' ')+'</button>';
});
h+='</div>';

// Custom field display
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var v=m[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="card-extra-row">';
customRows+='<span class="card-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="card-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="card-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

async function quickSet(id,availability){
try{
await fetch(A+'/members/'+id+'/availability',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({availability:availability})});
load();
}catch(e){alert('Update failed')}
}

// ─── Modal ────────────────────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=String(o).charAt(0).toUpperCase()+String(o).slice(1).replace(/_/g,' ');
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="2"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='datetime-local'){
var local='';
if(v){try{var d=new Date(v);if(!isNaN(d.getTime())){var pad=function(n){return n<10?'0'+n:''+n};local=d.getFullYear()+'-'+pad(d.getMonth()+1)+'-'+pad(d.getDate())+'T'+pad(d.getHours())+':'+pad(d.getMinutes())}}catch(e){}}
h+='<input type="datetime-local" id="f-'+f.name+'" value="'+esc(local)+'"'+ph+'>';
}else if(f.type==='number'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(member){
var m=member||{};
var isEdit=!!member;
var h='<h2>'+(isEdit?'EDIT MEMBER':'NEW MEMBER')+'</h2>';

h+=fieldHTML(fieldByName('name'),m.name);
h+='<div class="row2">'+fieldHTML(fieldByName('email'),m.email)+fieldHTML(fieldByName('role'),m.role)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('availability'),m.availability||'available')+fieldHTML(fieldByName('timezone'),m.timezone)+'</div>';
h+=fieldHTML(fieldByName('status_message'),m.status_message);
h+=fieldHTML(fieldByName('until_time'),m.until_time);

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,m[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit)h+='<button class="btn btn-del" onclick="delMember()">Delete</button>';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button>';
h+='</div>';
return h;
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var m=null;
for(var i=0;i<members.length;i++){if(members[i].id===id){m=members[i];break}}
if(!m)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(m);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Name is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='number')val=parseFloat(el.value)||0;
else if(f.type==='datetime-local'){
val='';
if(el.value){try{val=new Date(el.value).toISOString()}catch(e){}}
}else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/members/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/members',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Add failed');return}
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

async function delMember(){
if(!editId)return;
if(!confirm('Delete this member?'))return;
await fetch(A+'/members/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// Auto-refresh every 30s to pick up auto-cleared availability
setInterval(load,30000);

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.availability_states)){
var avField=fieldByName('availability');
if(avField)avField.options=cfg.availability_states;
}

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
