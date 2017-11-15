import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { environment as env } from '../../../environments/environment';
import { ApiService, NodeServices, App, Transports, NodeInfo, Message, FeedBackItem } from '../../service';
import { MatSnackBar, MatDialog, MatDialogRef } from '@angular/material';
import { DataSource } from '@angular/cdk/collections';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { Subject } from 'rxjs/Subject';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { UpdateCardComponent } from '../../components/update-card/update-card.component';
import { AlertComponent } from '../../components/alert/alert.component';
import 'rxjs/add/observable/of';
import 'rxjs/add/observable/timer';
import 'rxjs/add/operator/debounceTime';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Component({
  selector: 'app-sub-status',
  templateUrl: './sub-status.component.html',
  styleUrls: ['./sub-status.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class SubStatusComponent implements OnInit, OnDestroy {
  task = new Subject();
  alertMsg = '';
  sshColumns = ['index', 'key', 'del'];
  displayedColumns = ['index', 'key', 'app'];
  transportColumns = ['index', 'fromApp', 'fromNode', 'toNode', 'toApp'];
  appSource: SubStatusDataSource = null;
  sshSource: SubStatusDataSource = null;
  sockSource: SubStatusDataSource = null;
  transportSource: SubStatusDataSource = null;
  key = '';
  power = '';
  transports: Array<Transports> = [];
  status: NodeServices = null;
  apps: Array<App> = [];
  isManager = false;
  socketColor = 'close-status';
  sshColor = 'close-status';
  socketClientColor = 'close-status';
  sshClientColor = 'close-status';
  statrStatusCss = 'mat-primary';
  dialogTitle = '';
  sshTextarea = '';
  sshAllowNodes = [];
  sockTextarea = '';
  sockAllowNodes = [];
  sshConnectKey = '';
  taskTime = 1000;
  timer: Subscription = null;
  startRequest = false;
  feedBacks: Array<FeedBackItem> = [];
  sshClientForm = new FormGroup({
    nodeKey: new FormControl('', Validators.required),
    appKey: new FormControl('', Validators.required),
  });
  socketClientForm = new FormGroup({
    nodeKey: new FormControl('', Validators.required),
    appKey: new FormControl('', Validators.required),
  });
  sshClientPort = 0;
  socketClientPort = 0;
  nodeVersion = '';
  nodeTag = '';
  _appData = new SubDatabase();
  _transportData = new SubDatabase();
  _sshServerData = new SubDatabase();
  _socketServerData = new SubDatabase();
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private api: ApiService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog) {
  }

  ngOnInit() {
    this.appSource = new SubStatusDataSource(this._appData);
    this.transportSource = new SubStatusDataSource(this._transportData);
    this.sshSource = new SubStatusDataSource(this._sshServerData);
    this.sockSource = new SubStatusDataSource(this._socketServerData);
    if (env.taskTime) {
      this.taskTime = env.taskTime;
    }
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
  connectSocket(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    if (this.socketClientForm.valid) {
      const data = new FormData();
      data.append('toNode', this.socketClientForm.get('nodeKey').value);
      data.append('toApp', this.socketClientForm.get('appKey').value);
      console.log('data:', this.status.addr);
      this.api.connectSocketClicent(this.status.addr, data).subscribe(result => {
        console.log('conect socket client');
        this.task.next();
      });
    }
    this.dialog.closeAll();
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
        console.log('conect ssh client');
        this.task.next();
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
        this.task.next();
      }
    });
  }
  delAllowSockNode(ev: Event, index: number) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.sockAllowNodes.splice(index, 1);
    const data = new FormData();
    data.append('data', this.sockAllowNodes.toString());
    this.api.runSockServer(this.status.addr, data).subscribe(result => {
      if (result) {
        console.log('set socks result:', result);
        this.sockTextarea = '';
        this.task.next();
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
        this.task.next();
      }
    });
  }
  setSock(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    let dataStr = '';
    if (this.sockAllowNodes.length > 0 && this.sockTextarea.trim()) {
      dataStr = this.sockAllowNodes + ',' + this.sockTextarea.trim();
    } else {
      dataStr = this.sockAllowNodes.toString();
    }
    const data = new FormData();
    data.append('data', dataStr);
    this.api.runSSHServer(this.status.addr, data).subscribe(result => {
      if (result) {
        console.log('set ssh result:', result);
        this.sshTextarea = '';
        this.task.next();
      }
    });
  }
  checkUpdate(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const ref = this.dialog.open(UpdateCardComponent, {
      panelClass: 'update-panel',
      width: '370px',
      disableClose: true,
      data: {
        vesrion: this.nodeVersion,
        tag: this.nodeTag
      }
    });
    ref.componentInstance.nodeUrl = this.status.addr;
  }
  refresh(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.task.next();
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
        console.log('start socket server');
        this.task.next();
      }
    });
  }
  runSSHServer(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.api.runSSHServer(this.status.addr).subscribe(isOk => {
      if (isOk) {
        console.log('start ssh server');
        this.task.next();
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
      this._sshServerData.push(nodes);
    }
    this.dialogTitle = title;
    this.dialog.open(content, {
      width: '800px',
    });
  }
  openSockSettings(ev: Event, content: any, title: string) {
    if (!title) {
      return;
    }
    let nodes = [];
    if (this.status.apps && this.findService('sshs')) {
      console.log('get socks nodes');
      nodes = this.findService('socks').allow_nodes;
      this._socketServerData.push(nodes);
    }
    this.dialogTitle = title;
    this.dialog.open(content, {
      width: '800px',
    });
  }
  findService(search: string): App {
    if (!this.status || !this.status.apps) {
      return null;
    }
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
    this.task.debounceTime(1000).subscribe(() => {
      this.init();
    });
    this.timer = Observable.timer(0, this.taskTime).subscribe(() => {
      this.task.next();
    });
  }
  close() {
    this.timer.unsubscribe();
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
        this.sockAllowNodes = this.findService('socks').allow_nodes;
      }
      this._sshServerData.push(this.sshAllowNodes);
      this._socketServerData.push(this.sockAllowNodes);
      if (this.isExist('sshs')) {
        this.sshColor = this.statrStatusCss;
      }
      if (this.isExist('sshc')) {
        this.sshClientColor = this.statrStatusCss;
      }
      if (this.isExist('socks')) {
        this.socketColor = this.statrStatusCss;
      }
      if (this.isExist('sockc')) {
        this.socketClientColor = this.statrStatusCss;
      }
    }
  }
  fillTransport() {
    if (env.isManager && this.status.addr) {
      this.transports = [];
      this.api.getNodeInfo(this.status.addr).subscribe((info: NodeInfo) => {
        if (info) {
          this.nodeVersion = info.version;
          this.nodeTag = info.tag;
          this.transports = info.transports;
          this._transportData.push(info.transports);
          this.feedBacks = info.app_feedbacks;
          this.showMessage(info.messages);
        }
      }, err => {
        this._transportData.push(null);
      });
    }
  }
  getClientPort(client: string, target: any) {
    const app = this.findService(client);
    if (!app || !this.feedBacks) {
      return;
    } else {
      const result = this.feedBacks.find(el => {
        return el.key === app.key;
      });
      target = result.feedbacks.port ? result.feedbacks.port : 0;
    }
  }
  showMessage(msgs: Array<Array<Message>>) {
    if (!msgs || msgs[0] == null) {
      return;
    } else {
      msgs.sort((m1, m2) => {
        m1.sort(this.compareMsg);
        m2.sort(this.compareMsg);
        if (m1[0].priority < m2[0].priority) {
          return 1;
        }
        if (m1[0].priority < m2[0].priority) {
          return -1;
        }
        return 0;
      });
      this.alertMsg = msgs[0][0].msg;
      setTimeout(() => {
        this.dialog.open(AlertComponent, {
          width: '45rem',
          panelClass: 'alert',
          data: {
            msg: msgs[0][0].msg
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
    if (env.isManager) {
      this.api.getApps(this.status.addr).subscribe((apps: Array<App>) => {
        this.status.apps = apps;
        if (apps) {
          this.setServiceStatus();
        }
        this._appData.push(this.status.apps);
        this.getClientPort('sshc', this.sshClientPort);
        this.getClientPort('sockc', this.sshClientPort);
      }, err => {
        this._appData.push(null);
      });
    }
  }
  init() {
    this.startRequest = true;
    const data = new FormData();
    data.append('key', this.key);
    this.api.getNodeStatus(data).subscribe((nodeServices: NodeServices) => {
      if (nodeServices) {
        this.status = nodeServices;
        this.fillTransport();
        this.fillApps();
        // Set SSH Client Port

      }
    });
    // if (!this.startRequest) {

    // }
  }
}


export class SubDatabase {
  /** Stream that emits whenever the data has been modified. */
  dataChange: BehaviorSubject<any> = new BehaviorSubject<any>([]);
  isShow = false;
  get data(): any[] { return this.dataChange.value; }

  constructor() { }

  push(value: Array<any>) {
    if (!value || value.length <= 0) {
      this.isShow = false;
    } else {
      this.isShow = true;
    }
    this.dataChange.next(value);
  }
}

export class SubStatusDataSource extends DataSource<any> {
  size = 0;
  constructor(private _data: SubDatabase) {
    super();
  }
  connect(): Observable<any> {
    return this._data.dataChange;
  }

  disconnect() {
  }
}
