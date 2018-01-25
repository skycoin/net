import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ApiService, ConnectServiceInfo } from '../../service';
import { MatDialogRef } from '@angular/material';

@Component({
  // tslint:disable-next-line:component-selector
  selector: 'apps-setting',
  templateUrl: 'apps-setting.component.html',
  styleUrls: ['./apps-setting.component.scss']
})

export class AppsSettingComponent implements OnInit {
  socks_items: Array<AppSetting> = [];
  ssh_items: Array<AppSetting> = [];
  settingForm = new FormGroup({
    sshs: new FormControl('', Validators.required),
    sshc: new FormControl('', Validators.required),
    sshc_conf: new FormControl('', Validators.required),
    sshc_conf_nodeKey: new FormControl(''),
    sshc_conf_appKey: new FormControl(''),
    sockss: new FormControl('', Validators.required),
    socksc: new FormControl('', Validators.required),
    socksc_conf: new FormControl('', Validators.required),
    socksc_conf_nodeKey: new FormControl(''),
    socksc_conf_appKey: new FormControl(''),
  });
  sshc = 'sshc';
  sshs = 'sshs';
  sockss = 'sockss';
  socksc = 'socksc';
  socksc_opts: Array<ConnectServiceInfo> = [];
  sshc_opts: Array<ConnectServiceInfo> = [];
  addr = '';
  constructor(private api: ApiService, private dialogRef: MatDialogRef<AppsSettingComponent>) { }

  ngOnInit() {
    const data = new FormData();
    data.append('client', this.socksc);
    this.api.getClientConnection(data).subscribe(info => {
      this.socksc_opts = info;
      data.set('client', this.sshc);
      this.api.getClientConnection(data).subscribe(i => {
        this.sshc_opts = i;
      });
    });

    this.api.getAutoStart(this.addr).subscribe((config) => {
      this.settingForm.get('sockss').setValue(config.sockss === 'true' ? true : false);
      this.settingForm.get('socksc').setValue(config.socksc === 'true' ? true : false);
      this.settingForm.get('sshs').setValue(config.sshs === 'true' ? true : false);
      this.settingForm.get('sshc').setValue(config.sshc === 'true' ? true : false);
      this.settingForm.get('sshc_conf_nodeKey').setValue(config.sshc_conf_nodeKey);
      this.settingForm.get('sshc_conf_appKey').setValue(config.sshc_conf_appKey);
      this.settingForm.get('socksc_conf_nodeKey').setValue(config.socksc_conf_nodeKey);
      this.settingForm.get('socksc_conf_appKey').setValue(config.socksc_conf_appKey);
    });
  }
  save() {
    const data = new FormData();
    const json = this.settingForm.value;
    json[this.sshc] = String(json[this.sshc]);
    json[this.sshs] = String(json[this.sshs]);
    json[this.sockss] = String(json[this.sockss]);
    json[this.socksc] = String(json[this.socksc]);
    json['sshc_conf'] = '';
    json['socksc_conf'] = '';
    data.append('data', JSON.stringify(json));
    this.api.setAutoStart(this.addr, data).subscribe((result) => {
      if (result) {
        this.dialogRef.close();
      }
    });
  }
  setKey(client: string) {
    const info: ConnectServiceInfo = this.settingForm.get(`${client}_conf`).value;
    if (info) {
      this.settingForm.get(`${client}_conf_nodeKey`).patchValue(info.nodeKey);
      this.settingForm.get(`${client}_conf_appKey`).patchValue(info.appKey);
    }
  }
}

export interface AppSetting {
  text?: string;
  type?: string;
  opts?: Array<string>;
}
