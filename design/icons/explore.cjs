const fs=require('fs'),path=require('path');
const {Resvg}=require(path.join(process.cwd(),'.cache/icon-tools/node_modules/@resvg/resvg-js'));
const dir=__dirname;fs.mkdirSync(path.join(dir,'drafts'),{recursive:true});
const jade='#173F36',ivory='#F5EDDA',red='#CB493B',gold='#C8A56B';
const rect=(x,y,w,h,r,fill)=>`<rect x="${x}" y="${y}" width="${w}" height="${h}" rx="${r}" fill="${fill}"/>`;
const circle=(x,y,r,fill)=>`<circle cx="${x}" cy="${y}" r="${r}" fill="${fill}"/>`;
const line=(d,color=ivory,width=26)=>`<path d="${d}" fill="none" stroke="${color}" stroke-width="${width}" stroke-linecap="round" stroke-linejoin="round"/>`;
const shape=(d,fill=ivory)=>`<path d="${d}" fill="${fill}"/>`;
const concepts=[
 ['Signal tile',rect(145,101,222,310,34,ivory)+circle(318,148,31,red)+line('M195 320h120',jade,24)],
 ['Folded lesson',shape('M116 139Q180 118 256 157Q332 118 396 139V360Q328 343 256 383Q184 343 116 360Z')+line('M256 163v209',jade,18)+circle(324,224,24,red)],
 ['Dragon ribbon',shape('M155 112H337L378 166L306 226H373V301L278 400H138L217 311H147V240L244 176H155Z',red)+shape('M205 258H300L245 307H205Z',ivory)],
 ['Bamboo lesson',line('M256 397V148',ivory,32)+shape('M250 252Q151 252 142 171Q231 176 250 252')+shape('M263 194Q359 194 370 112Q283 117 263 194')+line('M226 308h60',red,20)],
 ['The missing tile',rect(97,168,91,217,20,ivory)+rect(211,168,91,217,20,ivory)+rect(325,108,91,217,20,red)],
 ['Five blocks',rect(117,117,117,117,22,ivory)+rect(259,117,117,117,22,ivory)+rect(117,259,117,117,22,ivory)+rect(259,259,50,117,16,red)+rect(326,259,50,117,16,red)],
 ['Wind rose',shape('M256 92L296 215L420 256L296 297L256 420L216 297L92 256L216 215Z')+circle(256,256,35,red)],
 ['Riichi deposit',`<g transform="rotate(-35 256 256)">${rect(98,217,316,78,39,ivory)+circle(256,256,23,red)}</g>`],
 ['Tile fan',`<g transform="rotate(-23 256 367)">${rect(146,124,179,261,24,gold)}</g><g transform="rotate(4 256 367)">${rect(181,115,179,261,24,ivory)}</g><g transform="rotate(27 256 367)">${rect(215,126,179,261,24,red)}</g>`],
 ['Enso wait',line('M379 183A151 151 0 1 0 393 300',ivory,31)+rect(214,190,84,132,19,red)],
 ['Sparrow',shape('M123 301L221 208L304 237L368 164L390 190L354 257L301 332L200 347Z')+shape('M173 310L112 370L243 339',red)+circle(344,210,8,jade)],
 ['Pinfu weave',line('M122 194h268M122 318h268',ivory,36)+line('M192 121v270M320 121v270',ivory,36)+rect(238,238,36,36,8,red)],
 ['Tile conversation',shape('M132 116H380V329H265L182 400V329H132Z')+rect(199,166,114,112,22,jade)+circle(256,222,25,red)],
 ['Study light',shape('M256 131L116 360H396Z',gold)+rect(209,241,94,143,17,ivory)+rect(188,109,136,37,18,red)],
 ['Open tile',shape('M135 117H377V395H286V202H226V395H135Z')+circle(256,143,23,red)],
 ['Single circle',circle(256,256,147,ivory)+circle(256,256,101,jade)+circle(256,256,55,red)],
 ['Paper crane',shape('M89 175L251 241L423 121L323 302L254 342L191 310Z')+shape('M191 310L129 397L254 342',red)+shape('M251 241L323 302L254 342',gold)],
 ['Table compass',shape('M256 92L420 256L256 420L92 256Z')+shape('M256 151L361 256L256 361L151 256Z',jade)+rect(229,217,54,78,11,red)],
 ['Riichi ligature',shape('M154 112H285Q374 112 374 204Q374 264 312 284L384 400H298L233 289H225V400H154Z')+shape('M225 177H277Q303 177 303 204Q303 231 277 231H225Z',jade)+circle(155,113,23,red)],
 ['Ready steps',rect(112,288,83,108,17,ivory)+rect(215,221,83,175,17,ivory)+rect(318,155,83,241,17,ivory)+circle(153,182,33,red)]
];
function svg(mark,bg=jade){return `<svg xmlns="http://www.w3.org/2000/svg" width="1024" height="1024" viewBox="0 0 512 512">${rect(24,24,464,464,105,bg)}${mark}</svg>`}
const records=concepts.map(([name,mark],i)=>({id:String(i+1).padStart(2,'0'),name,svg:svg(mark)}));
if(require.main===module)for(const r of records){fs.writeFileSync(path.join(dir,'drafts',r.id+'.svg'),r.svg);fs.writeFileSync(path.join(dir,'drafts',r.id+'.png'),new Resvg(r.svg,{fitTo:{mode:'width',value:1024},font:{loadSystemFonts:false}}).render().asPng())}
function board(records,file){const W=1000,H=Math.ceil(records.length/5)*205+30;let s=`<svg xmlns="http://www.w3.org/2000/svg" width="${W}" height="${H}"><rect width="100%" height="100%" fill="#E8E6DF"/>`;records.forEach((r,i)=>{const x=i%5*200+30,y=Math.floor(i/5)*205+20;s+=`<image x="${x}" y="${y}" width="140" height="140" href="data:image/svg+xml;base64,${Buffer.from(r.svg).toString('base64')}"/><text x="${x+70}" y="${y+167}" text-anchor="middle" font-family="Segoe UI" font-size="14" fill="#233A32">${r.id}. ${r.name}</text>`});s+='</svg>';fs.writeFileSync(path.join(dir,file+'.svg'),s);fs.writeFileSync(path.join(dir,file+'.png'),new Resvg(s).render().asPng());}
if(require.main===module){board(records,'round-1');fs.writeFileSync(path.join(dir,'concepts.json'),JSON.stringify(records.map(({id,name})=>({id,name})),null,2));}
module.exports={svg,rect,circle,line,shape,board,jade,ivory,red,gold,dir,Resvg,records};
