package server
import "net/http"
func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html"); w.Write([]byte(dashHTML)) }
const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Semaphore</title><link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet"><style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}.main{padding:1.5rem;max-width:960px;margin:0 auto}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}.toolbar{display:flex;gap:.5rem;margin-bottom:1rem}.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.search:focus{outline:none;border-color:var(--leather)}.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}.item:hover{border-color:var(--leather)}.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}.item-title{font-size:.85rem;font-weight:700}.item-sub{font-size:.7rem;color:var(--cd);margin-top:.1rem}.item-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}.item-notes{font-size:.65rem;color:var(--cm);margin-top:.3rem;font-style:italic;padding:.3rem .5rem;border-left:2px solid var(--bg3);display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}.item-actions{display:flex;gap:.3rem;flex-shrink:0}.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm)}.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-sm{font-size:.55rem;padding:.2rem .4rem}.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> SEMAPHORE</h1><button class="btn btn-p" onclick="openForm()">+ Add</button></div>
<div class="main"><div class="stats" id="stats"></div><div class="toolbar"><input class="search" id="search" placeholder="Search..." oninput="render()"></div><div id="list"></div></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/team_members').then(function(r){return r.json()});items=r.team_members||[];document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+items.length+'</div><div class="st-l">Total</div></div><div class="st"><div class="st-v">-</div><div class="st-l">&nbsp;</div></div><div class="st"><div class="st-v">-</div><div class="st-l">&nbsp;</div></div>';render();}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(i){return (i.name||'').toLowerCase().includes(q)||(i.email||'').toLowerCase().includes(q)||(i.role||'').toLowerCase().includes(q)||(i.availability||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No items found.</div>';return;}
var h='';f.forEach(function(i){
h+='<div class="item"><div class="item-top"><div style="flex:1"><div class="item-title">'+esc(i.name)+'</div>';
if(i.email)h+='<div class="item-sub">'+esc(i.email)+'</div>';
h+='</div><div class="item-actions"><button class="btn btn-sm" onclick="openEdit(\''+i.id+'\')">Edit</button><button class="btn btn-sm" onclick="del(\''+i.id+'\')" style="color:var(--red)">&#10005;</button></div></div>';
h+='<div class="item-meta">';
if(i.role)h+='<span>'+esc(i.role)+'</span>';
if(i.availability)h+='<span>'+esc(i.availability)+'</span>';
if(i.status_message)h+='<span>'+esc(i.status_message)+'</span>';
if(i.timezone)h+='<span>'+esc(i.timezone)+'</span>';
if(i.until_time)h+='<span>'+esc(i.until_time)+'</span>';
h+='<span>'+ft(i.created_at)+'</span></div>';

h+='</div>';});document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/team_members/'+id,{method:'DELETE'});load();}
function formHTML(item){var i=item||{name:"",email:"",role:"",availability:"",status_message:"",timezone:"",until_time:""};var isEdit=!!item;
var h='<h2>'+(isEdit?'EDIT':'NEW')+' SEMAPHORE</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'"></div>';
h+='<div class="fr"><label>Email</label><input id="f-email" value="'+esc(i.email)+'"></div>';
h+='<div class="fr"><label>Role</label><input id="f-role" value="'+esc(i.role)+'"></div>';
h+='<div class="fr"><label>Availability</label><input id="f-availability" value="'+esc(i.availability)+'"></div>';
h+='<div class="fr"><label>Status Message</label><input id="f-status-message" value="'+esc(i.status_message)+'"></div>';
h+='<div class="fr"><label>Timezone</label><input id="f-timezone" value="'+esc(i.timezone)+'"></div>';
h+='<div class="fr"><label>Until Time</label><input id="f-until-time" value="'+esc(i.until_time)+'"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var x=null;for(var j=0;j<items.length;j++){if(items[j].id===id){x=items[j];break;}}if(!x)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(x);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var body={name:document.getElementById('f-name').value.trim(),email:document.getElementById('f-email').value.trim(),role:document.getElementById('f-role').value.trim(),availability:document.getElementById('f-availability').value.trim(),status_message:document.getElementById('f-status-message').value.trim(),timezone:document.getElementById('f-timezone').value.trim(),until_time:document.getElementById('f-until-time').value.trim()};
if(editId){await fetch(A+'/team_members/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/team_members',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}closeModal();load();}
function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});load();
</script></body></html>`
