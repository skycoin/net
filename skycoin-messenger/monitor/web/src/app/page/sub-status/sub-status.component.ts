import { Component, OnInit, ViewEncapsulation, OnDestroy, ViewChild, TemplateRef } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { environment as env } from '../../../environments/environment';
import { ApiService, NodeServices, App, Transports, NodeInfo, Message } from '../../service';
import { MatSnackBar, MatDialog, MatDialogRef } from '@angular/material';
import { DataSource } from '@angular/cdk/collections';
import { Observable } from 'rxjs/Observable';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { UpdateCardComponent } from '../../components/update-card/update-card.component';
import { AlertComponent } from '../../components/alert/alert.component';
import 'rxjs/add/observable/of';

@Component({
  selector: 'app-sub-status',
  templateUrl: './sub-status.component.html',
  styleUrls: ['./sub-status.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class SubStatusComponent implements OnInit, OnDestroy {
  @ViewChild('alert') alertRef: TemplateRef<any>;
  alertMsg = '';
  sshColumns = ['index', 'key', 'del'];
  displayedColumns = ['index', 'key', 'app'];
  transportColumns = ['index', 'fromApp', 'fromNode', 'toNode', 'toApp'];
  appSource: SubStatusDataSource = null;
  sshSource: SubStatusDataSource = null;
  transportSource: SubStatusDataSource = null;
  key = '';
  power = '';
  transports: Array<Transports> = [];
  status: NodeServices = null;
  apps: Array<App> = [];
  task = null;
  isManager = false;
  socketColor = 'close-status';
  sshColor = 'close-status';
  socketClientColor = 'close-status';
  sshClientColor = 'close-status';
  statrStatusCss = 'mat-primary';
  dialogTitle = '';
  sshTextarea = '';
  sshAllowNodes = [];
  socksTextarea = '';
  sshConnectKey = '';
  sshClientForm = new FormGroup({
    nodeKey: new FormControl('', Validators.required),
    appKey: new FormControl('', Validators.required),
  });
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private api: ApiService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog) {
  }

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      this.key = params.key;
      this.startTask();
      this.power = 'warn';
      this.isManager = env.isManager;
    });
  }
  ngOnDestroy() {
    this.close();
  }
  transportsTrackBy(index, transport) {
    return transport ? transport.from_node : undefined;
  }
  appTrackBy(index, app) {
    return app ? app.key : undefined;
  }
  connectSSH(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (this.sshClientForm.valid) {
      const data = new FormData();
      data.append('toNode', this.sshClientForm.get('nodeKey').value);
      data.append('toApp', this.sshClientForm.get('appKey').value);
      this.api.connectSSHClient(this.status.addr, data).subscribe(result => {
        console.log('conect ssh client:', result);
        this.init();
      });
    }
    this.dialog.closeAll();
  }
  delAllowNode(ev: Event, index: number) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.sshAllowNodes.splice(index, 1);
    const data = new FormData();
    data.append('data', this.sshAllowNodes.toString());
    this.api.runSSHServer(this.status.addr, data).subscribe(result => {
      if (result) {
        console.log('set ssh result:', result);
        this.sshTextarea = '';
        this.init();
      }
    });
  }
  setSSH(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    let dataStr = '';
    if (this.sshAllowNodes.length > 0 && this.sshTextarea.trim()) {
      dataStr = this.sshAllowNodes + ',' + this.sshTextarea.trim();
    } else {
      dataStr = this.sshAllowNodes.toString();
    }
    const data = new FormData();
    data.append('data', dataStr);
    this.api.runSSHServer(this.status.addr, data).subscribe(result => {
      if (result) {
        console.log('set ssh result:', result);
        this.sshTextarea = '';
        this.init();
      }
    });
  }
  checkUpdate(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.dialog.open(UpdateCardComponent, {
      panelClass: 'update-panel',
      width: '370px',
      disableClose: true
    });
  }
  refresh(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.init();
    this.snackBar.open('Refreshed', 'Dismiss', {
      duration: 3000,
      verticalPosition: 'top',
      extraClasses: ['bg-success']
    });
  }
  runSocketServer(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.api.runSockServer(this.status.addr).subscribe(isOk => {
      if (isOk) {
        this.init();
        console.log('start socket server');
      }
    });
  }
  runSSHServer(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.api.runSSHServer(this.status.addr).subscribe(isOk => {
      if (isOk) {
        this.init();
        console.log('start ssh server');
      }
    });
  }
  reboot(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    console.log('reboot');
    this.api.reboot(this.status.addr).subscribe(isOk => {
      if (isOk) {
        if (this.task) {
          this.close();
        }
        this.startTask();
      }
    });
  }
  inputKeys(ev: Event, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.dialog.open(content, {
      width: '450px'
    });
  }
  openSettings(ev: Event, content: any, title: string) {
    if (!title) {
      return;
    }
    let nodes = [];
    if (this.status.apps && this.findService('sshs')) {
      console.log('get sshs nodes');
      nodes = this.findService('sshs').allow_nodes;
      this.sshSource = new SubStatusDataSource(nodes);
    }
    this.dialogTitle = title;
    this.dialog.open(content, {
      width: '800px',
    });
  }
  findService(search: string): App {
    const result = this.status.apps.find(el => {
      const tmp = el.attributes.find(attr => {
        return search === attr;
      });
      return search === tmp;
    });
    return result;
  }
  startTask() {
    this.init();
    this.task = setInterval(() => {
      this.init();
    }, 15000);
  }
  close() {
    clearInterval(this.task);
    this.task = null;
  }
  isExist(search: string) {
    const result = this.status.apps.find(el => {
      const tmp = el.attributes.find(attr => {
        return search === attr;
      });
      return search === tmp;
    });
    return result !== undefined && result !== null;
  }
  setServiceStatus() {
    this.sshColor = 'close-status';
    this.sshClientColor = 'close-status';
    this.socketColor = 'close-status';
    this.socketClientColor = 'close-status';
    if (this.status.apps) {
      if (this.findService('sshs')) {
        this.sshAllowNodes = this.findService('sshs').allow_nodes;
      }
      this.sshSource = new SubStatusDataSource(this.sshAllowNodes);
      if (this.isExist('sshs')) {
        this.sshColor = this.statrStatusCss;
      }
      if (this.isExist('sshc')) {
        this.sshClientColor = this.statrStatusCss;
      }
      if (this.isExist('sockss')) {
        this.socketColor = this.statrStatusCss;
      }
      if (this.isExist('socksc')) {
        this.socketClientColor = this.statrStatusCss;
      }
    }
  }
  fillTransport() {
    if (env.isManager && this.status.addr) {
      // this.transportSource = null;
      this.transports = [];
      this.api.getNodeInfo(this.status.addr).subscribe((info: NodeInfo) => {
        if (info) {
          if (info.transports && info.transports.length > 0) {
            this.transports = info.transports;
            this.transportSource = new SubStatusDataSource(info.transports);
          }
          this.showMessage(info.messages);
        }
      });
    }
  }
  showMessage(msgs: Array<Array<Message>>) {
    if (!msgs || msgs[0] == null) {
      return;
    } else {
      msgs.sort((m1, m2) => {
        m1.sort(this.compareMsg);
        m2.sort(this.compareMsg);
        if (m1[0].Priority < m2[0].Priority) {
          return 1;
        }
        if (m1[0].Priority < m2[0].Priority) {
          return -1;
        }
        return 0;
      });
      this.alertMsg = msgs[0][0].Msg;
      this.dialog.open(this.alertRef, {
        width: '45rem',
        panelClass: 'alert'
      });
      setTimeout(() => {
        this.dialog.open(AlertComponent, {
          width: '45rem',
          panelClass: 'alert',
          data: {
            msg: msgs[0][0].Msg
          }
        });
      }, 500);
    }
  }
  compareMsg(msg1, msg2) {
    if (msg1.Priority < msg2.Priority) {
      return 1;
    }
    if (msg1.Priority > msg2.Priority) {
      return -1;
    }
    return 0;
  }
  fillApps() {
    if (this.status.apps) {
      this.appSource = new SubStatusDataSource(this.status.apps);
      this.setServiceStatus();
    } else {
      if (env.isManager) {
        this.api.getApps(this.status.addr).subscribe((apps: Array<App>) => {
          this.status.apps = apps;
          this.appSource = new SubStatusDataSource(this.status.apps);
          this.setServiceStatus();
        });
      }
    }
  }
  init() {
    const data = new FormData();
    data.append('key', this.key);
    this.api.getNodeStatus(data).subscribe((nodeServices: NodeServices) => {
      if (nodeServices) {
        this.status = nodeServices;
        this.fillTransport();
        this.fillApps();
      }
    });
  }
}
export class SubStatusDataSource extends DataSource<any> {
  size = 0;
  constructor(private apps: any) {
    super();
  }
  connect(): Observable<any> {
    return Observable.of(this.apps);
  }

  disconnect() { }
}
