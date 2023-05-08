// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {midihandler} from '../models';
import {gui} from '../models';

export function ClearCommand(arg1:midihandler.Event,arg2:boolean):Promise<void>;

export function GetActions():Promise<midihandler.Actions>;

export function GetMIDIDevices():Promise<Array<string>>;

export function ListenForInput():Promise<midihandler.Event>;

export function LoadSettings():Promise<gui.Settings>;

export function SelectDevice(arg1:string):Promise<string>;

export function SetCommand(arg1:midihandler.Event,arg2:boolean,arg3:string):Promise<void>;
