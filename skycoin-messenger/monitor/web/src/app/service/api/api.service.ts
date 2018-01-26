import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';
import 'rxjs/add/observable/empty';
import { Router } from '@angular/router';
import { MatDialog } from '@angular/material';
import { AlertService } from '../alert/alert.service';

@Injectable()
export class ApiService {
  private connUrl = '/conn/';
  private nodeUrl = '/node';
  private bankUrl = '127.0.0.1:8080/';
  private callbackParm = 'callback';
  private jsonHeader = { 'Content-Type': 'application/json' };
  constructor(
    private httpClient: HttpClient,
    private router: Router,
    private dialog: MatDialog,
    private alert: AlertService) { }

  getServerInfo() {
    return this.handleGet(this.connUrl + 'getServerInfo', { responseType: 'text' });
  }
  closeApp(addr: string, data: FormData) {
    return this.handleNodePost(addr, '/node/run/closeApp', data);
  }
  login(data: FormData) {
    return this.handlePost('login', data);
  }
  updatePass(data: FormData) {
    return this.handlePost('updatePass', data);
  }
  checkLogin() {
    return this.handlePost('checkLogin', null, { responseType: 'text' });
  }
  getAllNode() {
    return this.handleGet(this.connUrl + 'getAll');
  }
  getNodeStatus(data: FormData) {
    return this.handlePost(this.connUrl + 'getNode', data);
  }
  setNodeConfig(data: FormData) {
    return this.handlePost(this.connUrl + 'setNodeConfig', data);
  }
  updateNodeConfig(addr: string) {
    return this.handleNodePost(addr, '/node/run/updateNode');
  }
  getMsgs(addr) {
    return this.handleNodePost(addr, '/node/getMsgs');
  }
  getApps(addr: string) {
    return this.handleNodePost(addr, '/node/getApps');
  }

  getNodeInfo(addr: string) {
    return this.handleNodePost(addr, '/node/getInfo');
  }
  reboot(addr: string) {
    return this.handleNodePost(addr, '/node/reboot');
  }
  getAutoStart(addr: string, data: FormData) {
    return this.handleNodePost(addr, '/node/run/getAutoStartConfig', data);
  }
  setAutoStart(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/run/setAutoStartConfig', data);
  }
  checkAppMsg(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/getMsg', data);
  }
  searchServices(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/run/searchServices', data);
  }
  getServicesResult(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/run/getSearchServicesResult', data);
  }
  connectSSHClient(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/run/sshc', data);
  }
  connectSocketClicent(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/run/socksc', data);
  }
  runSSHServer(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/run/sshs', data);
  }
  runSockServer(addr: string, data?: FormData) {
    return this.handleNodePost(addr, '/node/run/sockss', data);
  }
  runNodeupdate(addr: string) {
    return this.handleNodePost(addr, '/node/run/update');
  }
  getDebugPage(addr: string) {
    return this.handleNodePost(addr, '/debug/pprof');
  }
  checkUpdate(channel, vesrion: string) {
    const data = new FormData();
    data.append('addr', `http://messenger.skycoin.net:8100/api/version?c=${channel}&v=${vesrion}`);
    return this.handlePost(this.nodeUrl, data);
  }
  saveClientConnection(data: FormData) {
    return this.handlePost(this.connUrl + 'saveClientConnection', data);
  }
  removeClientConnection(data: FormData) {
    return this.handlePost(this.connUrl + 'removeClientConnection', data);
  }
  editClientConnection(data: FormData) {
    return this.handlePost(this.connUrl + 'editClientConnection', data);
  }
  SetClientAutoStart(data: FormData) {
    return this.handlePost(this.connUrl + 'setClientAutoStart', data);
  }
  getClientConnection(data: FormData) {
    return this.handlePost(this.connUrl + 'getClientConnection', data);
  }

  getWalletNewAddress(data: FormData) {
    return this.handleNodePost(this.bankUrl, 'skypay/tools/newAddress', data);
  }
  getWalletInfo(data: FormData) {
    return this.handleNodePost(this.bankUrl, 'skypay/node/get', data);
  }
  jsonp(url: string) {
    if (url === '') {
      return Observable.throw('Url is empty.');
    }
    return this.httpClient.jsonp(url, this.callbackParm).catch(err => Observable.throw(err));
  }
  handleGet(url: string, opts?: any) {
    if (url === '') {
      return Observable.throw('Url is empty.');
    }
    return this.httpClient.get(url, opts).catch(err => this.handleError(err));
  }
  handleNodePost(nodeAddr: string, api: string, data?: FormData, opts?: any) {
    if (nodeAddr === '' || api === '') {
      return Observable.throw('nodeAddr or api is empty.');
    }
    nodeAddr = 'http://' + nodeAddr + api;
    if (!data) {
      data = new FormData();
    }
    data.append('addr', nodeAddr);
    return this.handlePost(this.nodeUrl, data, opts);
  }
  handlePost(url: string, data?: FormData, opts?: any) {
    if (url === '') {
      return Observable.throw('Url is empty.');
    }
    return this.httpClient.post(url, data, opts).catch(err => this.handleNodeError(err));
  }
  handleNodeError(err: HttpErrorResponse) {
    if (err.status === 302) {
      console.log('open url');
      this.dialog.closeAll();
      this.router.navigate([{ outlets: { user: ['login'] } }]);
      return Observable.empty();
    }
    return Observable.throw(err.error.text);
  }
  handleError(err: HttpErrorResponse) {
    if (err.status === 302) {
      console.log('open url');
      this.dialog.closeAll();
      this.router.navigate([{ outlets: { user: ['login'] } }]);
      return Observable.empty();
    }
    return Observable.throw(err);
  }
}
export interface ConnectServiceInfo {
  label?: string;
  nodeKey?: string;
  appKey?: string;
  count?: number;
  auto_start?: boolean;
}
export interface Conn {
  key?: string;
  type?: string;
  send_bytes?: number;
  recv_bytes?: number;
  last_ack_time?: number;
  start_time?: number;
}
export interface ConnData extends Conn {
  index?: number;
}
export interface ConnsResponse {
  conns?: Array<Conn>;
}

export interface NodeServices extends Conn {
  apps?: Array<App>;
  addr?: string;
}

export interface App {
  key?: string;
  attributes?: Array<string>;
  allow_nodes?: Array<string>;
}

export interface Transports {
  from_node?: string;
  to_node?: string;
  from_app?: string;
  to_app?: string;
}
// export interface Message {
//   priority?: number;
//   type?: number;
//   msg?: string;
// }
export interface FeedBack {
  port?: number;
  failed?: boolean;
  msg?: Message;
}
export interface FeedBackItem {
  key?: string;
  failed?: boolean;
  port?: number;
  unread?: boolean;
}
export interface NodeInfo {
  version?: string;
  tag?: string;
  discoveries?: Map<string, boolean>;
  transports?: Array<Transports>;
  messages?: Array<Message>;
  app_feedbacks?: Array<FeedBackItem>;
}

export interface Message {
  key?: string;
  read?: boolean;
  msgs?: Array<MessageItem>;
}
export interface MessageItem {
  msg?: string;
  priority?: number;
  time?: number;
  type?: number;
}

export interface AutoStartConfig {
  socks_server?: boolean;
  ssh_server?: boolean;
}

export interface WalletAddress {
  code?: number;
  address?: string;
}
