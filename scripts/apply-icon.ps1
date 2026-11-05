param(
    [Parameter(Mandatory=$true)][string]$Executable,
    [string]$Icon = "$PSScriptRoot/../resources/icon.ico"
)
$ErrorActionPreference = 'Stop'
$exePath = (Resolve-Path -LiteralPath $Executable).Path
$iconPath = (Resolve-Path -LiteralPath $Icon).Path
if ([IO.Path]::GetExtension($exePath) -ne '.exe') { throw 'Expected a Windows executable.' }
# Preserve MyGo's manifest/version resources while installing the optically
# adjusted PNG frames of our ICO, rather than automatic downscales of the master.
Add-Type -TypeDefinition @'
using System;
using System.ComponentModel;
using System.IO;
using System.Runtime.InteropServices;
public static class RiichiIconResources {
 [DllImport("kernel32.dll",CharSet=CharSet.Unicode,SetLastError=true)]
 static extern IntPtr BeginUpdateResource(string file,bool deleteExisting);
 [DllImport("kernel32.dll",SetLastError=true)]
 static extern bool UpdateResource(IntPtr handle,IntPtr type,IntPtr name,ushort language,byte[] data,uint size);
 [DllImport("kernel32.dll",SetLastError=true)]
 static extern bool EndUpdateResource(IntPtr handle,bool discard);
 public static void Apply(string exe,string ico) {
  byte[] bytes=File.ReadAllBytes(ico);
  using(var input=new BinaryReader(new MemoryStream(bytes))) {
   if(input.ReadUInt16()!=0 || input.ReadUInt16()!=1)throw new Exception("Invalid ICO");
   int count=input.ReadUInt16();
   IntPtr handle=BeginUpdateResource(exe,false);
   if(handle==IntPtr.Zero)throw new Win32Exception();
   bool committed=false;
   try {
    using(var groupStream=new MemoryStream())using(var group=new BinaryWriter(groupStream)) {
     group.Write((ushort)0);group.Write((ushort)1);group.Write((ushort)count);
     for(int i=0;i<count;i++) {
      input.BaseStream.Position=6+i*16;
      byte[] head=input.ReadBytes(8);uint length=input.ReadUInt32();uint offset=input.ReadUInt32();
      input.BaseStream.Position=offset;byte[] png=input.ReadBytes((int)length);
      if(png.Length!=length)throw new Exception("Truncated ICO");
      if(!UpdateResource(handle,(IntPtr)3,(IntPtr)(i+1),1033,png,length))throw new Win32Exception();
      group.Write(head);group.Write(length);group.Write((ushort)(i+1));
     }
     byte[] data=groupStream.ToArray();
     if(!UpdateResource(handle,(IntPtr)14,(IntPtr)1,1033,data,(uint)data.Length))throw new Win32Exception();
    }
    if(!EndUpdateResource(handle,false))throw new Win32Exception();
    committed=true;
   } finally { if(!committed)EndUpdateResource(handle,true); }
  }
 }
}
'@
[RiichiIconResources]::Apply($exePath,$iconPath)
Write-Output "Installed the nine icon sizes in $exePath"
