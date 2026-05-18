export namespace analysis {
	
	export class BreakdownRow {
	    name: string;
	    matches: number;
	    rounds: number;
	    kd: number;
	    winRate: number;
	    headshotPercent: number;
	
	    static createFrom(source: any = {}) {
	        return new BreakdownRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.matches = source["matches"];
	        this.rounds = source["rounds"];
	        this.kd = source["kd"];
	        this.winRate = source["winRate"];
	        this.headshotPercent = source["headshotPercent"];
	    }
	}
	export class Evidence {
	    matchIds: string[];
	    map?: string;
	    agent?: string;
	    metric: string;
	    value: number;
	    sampleSize: number;
	    comparisonBaseline: number;
	
	    static createFrom(source: any = {}) {
	        return new Evidence(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matchIds = source["matchIds"];
	        this.map = source["map"];
	        this.agent = source["agent"];
	        this.metric = source["metric"];
	        this.value = source["value"];
	        this.sampleSize = source["sampleSize"];
	        this.comparisonBaseline = source["comparisonBaseline"];
	    }
	}
	export class Finding {
	    id: string;
	    title: string;
	    severity: string;
	    confidence: string;
	    detail: string;
	    evidence: Evidence[];
	
	    static createFrom(source: any = {}) {
	        return new Finding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.severity = source["severity"];
	        this.confidence = source["confidence"];
	        this.detail = source["detail"];
	        this.evidence = this.convertValues(source["evidence"], Evidence);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MatchSummary {
	    id: string;
	    map: string;
	    agent: string;
	    role: string;
	    kills: number;
	    deaths: number;
	    assists: number;
	    roundsPlayed: number;
	    firstBloods: number;
	    firstDeaths: number;
	    headshotPercent: number;
	    won: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MatchSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.map = source["map"];
	        this.agent = source["agent"];
	        this.role = source["role"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.assists = source["assists"];
	        this.roundsPlayed = source["roundsPlayed"];
	        this.firstBloods = source["firstBloods"];
	        this.firstDeaths = source["firstDeaths"];
	        this.headshotPercent = source["headshotPercent"];
	        this.won = source["won"];
	    }
	}
	export class MetricSummary {
	    matches: number;
	    rounds: number;
	    kd: number;
	    kda: number;
	    headshotPercent: number;
	    firstBloodRate: number;
	    firstDeathRate: number;
	    winRate: number;
	    weakestMap: string;
	    weakestMapWinRate: number;
	    weakestMapSample: number;
	    primaryRoleObserved: string;
	
	    static createFrom(source: any = {}) {
	        return new MetricSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matches = source["matches"];
	        this.rounds = source["rounds"];
	        this.kd = source["kd"];
	        this.kda = source["kda"];
	        this.headshotPercent = source["headshotPercent"];
	        this.firstBloodRate = source["firstBloodRate"];
	        this.firstDeathRate = source["firstDeathRate"];
	        this.winRate = source["winRate"];
	        this.weakestMap = source["weakestMap"];
	        this.weakestMapWinRate = source["weakestMapWinRate"];
	        this.weakestMapSample = source["weakestMapSample"];
	        this.primaryRoleObserved = source["primaryRoleObserved"];
	    }
	}
	export class PlayerSnapshot {
	    name: string;
	    tagline: string;
	    region: string;
	    primaryRole: string;
	    recentMatches: MatchSummary[];
	
	    static createFrom(source: any = {}) {
	        return new PlayerSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.tagline = source["tagline"];
	        this.region = source["region"];
	        this.primaryRole = source["primaryRole"];
	        this.recentMatches = this.convertValues(source["recentMatches"], MatchSummary);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PracticeTask {
	    day: number;
	    focus: string;
	    map?: string;
	    agent?: string;
	    duration: string;
	    checklist: string[];
	    evidence: string;
	
	    static createFrom(source: any = {}) {
	        return new PracticeTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.day = source["day"];
	        this.focus = source["focus"];
	        this.map = source["map"];
	        this.agent = source["agent"];
	        this.duration = source["duration"];
	        this.checklist = source["checklist"];
	        this.evidence = source["evidence"];
	    }
	}
	export class Recommendation {
	    id: string;
	    findingId: string;
	    title: string;
	    reason: string;
	    drill: string;
	    cadence: string;
	
	    static createFrom(source: any = {}) {
	        return new Recommendation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.findingId = source["findingId"];
	        this.title = source["title"];
	        this.reason = source["reason"];
	        this.drill = source["drill"];
	        this.cadence = source["cadence"];
	    }
	}
	export class Report {
	    player: PlayerSnapshot;
	    metrics: MetricSummary;
	    mapBreakdown: BreakdownRow[];
	    agentBreakdown: BreakdownRow[];
	    practicePlan: PracticeTask[];
	    findings: Finding[];
	    recommendations: Recommendation[];
	
	    static createFrom(source: any = {}) {
	        return new Report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player = this.convertValues(source["player"], PlayerSnapshot);
	        this.metrics = this.convertValues(source["metrics"], MetricSummary);
	        this.mapBreakdown = this.convertValues(source["mapBreakdown"], BreakdownRow);
	        this.agentBreakdown = this.convertValues(source["agentBreakdown"], BreakdownRow);
	        this.practicePlan = this.convertValues(source["practicePlan"], PracticeTask);
	        this.findings = this.convertValues(source["findings"], Finding);
	        this.recommendations = this.convertValues(source["recommendations"], Recommendation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace assistant {
	
	export class Alert {
	    id: string;
	    title: string;
	    message: string;
	    severity: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new Alert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.message = source["message"];
	        this.severity = source["severity"];
	        this.source = source["source"];
	    }
	}
	export class SessionState {
	    active: boolean;
	    startedAt: string;
	    roundCount: number;
	    tipsShown: number;
	    lastAlertAt: string;
	    currentAlert?: Alert;
	    message: string;
	    queueSize: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.active = source["active"];
	        this.startedAt = source["startedAt"];
	        this.roundCount = source["roundCount"];
	        this.tipsShown = source["tipsShown"];
	        this.lastAlertAt = source["lastAlertAt"];
	        this.currentAlert = this.convertValues(source["currentAlert"], Alert);
	        this.message = source["message"];
	        this.queueSize = source["queueSize"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TipResult {
	    hasTip: boolean;
	    alert: Alert;
	    state: SessionState;
	
	    static createFrom(source: any = {}) {
	        return new TipResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hasTip = source["hasTip"];
	        this.alert = this.convertValues(source["alert"], Alert);
	        this.state = this.convertValues(source["state"], SessionState);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace henrik {
	
	export class Status {
	    baseURL: string;
	    consentGranted: boolean;
	    canFetchPersonalData: boolean;
	    rateLimitPerMinute: number;
	    cacheTTLMinutes: number;
	    safeMode: boolean;
	    message: string;
	    nextStep: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.baseURL = source["baseURL"];
	        this.consentGranted = source["consentGranted"];
	        this.canFetchPersonalData = source["canFetchPersonalData"];
	        this.rateLimitPerMinute = source["rateLimitPerMinute"];
	        this.cacheTTLMinutes = source["cacheTTLMinutes"];
	        this.safeMode = source["safeMode"];
	        this.message = source["message"];
	        this.nextStep = source["nextStep"];
	    }
	}

}

export namespace localstore {
	
	export class PracticeProgressState {
	    items: Record<string, boolean>;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new PracticeProgressState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = source["items"];
	        this.updatedAt = source["updatedAt"];
	    }
	}

}

export namespace practice {
	
	export class Session {
	    id: string;
	    taskId: string;
	    focus: string;
	    map: string;
	    agent: string;
	    durationSeconds: number;
	    startedAt: string;
	    finishedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.taskId = source["taskId"];
	        this.focus = source["focus"];
	        this.map = source["map"];
	        this.agent = source["agent"];
	        this.durationSeconds = source["durationSeconds"];
	        this.startedAt = source["startedAt"];
	        this.finishedAt = source["finishedAt"];
	    }
	}
	export class SessionInput {
	    taskId: string;
	    focus: string;
	    map: string;
	    agent: string;
	    durationSeconds: number;
	    startedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.taskId = source["taskId"];
	        this.focus = source["focus"];
	        this.map = source["map"];
	        this.agent = source["agent"];
	        this.durationSeconds = source["durationSeconds"];
	        this.startedAt = source["startedAt"];
	    }
	}
	export class SessionState {
	    sessions: Session[];
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessions = this.convertValues(source["sessions"], Session);
	        this.updatedAt = source["updatedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace settings {
	
	export class DataSettings {
	    consentPersonalData: boolean;
	    riotName: string;
	    riotTag: string;
	    puuid: string;
	    region: string;
	    shard: string;
	    apiKey: string;
	    apiKeyHeader: string;
	    rateLimitTier: string;
	    matchCount: number;
	    cacheTTLMinutes: number;
	    lastUpdatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new DataSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.consentPersonalData = source["consentPersonalData"];
	        this.riotName = source["riotName"];
	        this.riotTag = source["riotTag"];
	        this.puuid = source["puuid"];
	        this.region = source["region"];
	        this.shard = source["shard"];
	        this.apiKey = source["apiKey"];
	        this.apiKeyHeader = source["apiKeyHeader"];
	        this.rateLimitTier = source["rateLimitTier"];
	        this.matchCount = source["matchCount"];
	        this.cacheTTLMinutes = source["cacheTTLMinutes"];
	        this.lastUpdatedAt = source["lastUpdatedAt"];
	    }
	}

}

export namespace tactical {
	
	export class MapCatalogEntry {
	    id: string;
	    uuid: string;
	    name: string;
	    displayName: string;
	    imageUrl: string;
	    tacticalImageUrl: string;
	    hasTacticalLayout: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MapCatalogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.uuid = source["uuid"];
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	        this.imageUrl = source["imageUrl"];
	        this.tacticalImageUrl = source["tacticalImageUrl"];
	        this.hasTacticalLayout = source["hasTacticalLayout"];
	    }
	}
	export class PlanLine {
	    id: string;
	    label: string;
	    x1: number;
	    y1: number;
	    x2: number;
	    y2: number;
	
	    static createFrom(source: any = {}) {
	        return new PlanLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.x1 = source["x1"];
	        this.y1 = source["y1"];
	        this.x2 = source["x2"];
	        this.y2 = source["y2"];
	    }
	}
	export class PlanMarker {
	    id: string;
	    kind: string;
	    label: string;
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new PlanMarker(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.kind = source["kind"];
	        this.label = source["label"];
	        this.x = source["x"];
	        this.y = source["y"];
	    }
	}
	export class MapPlan {
	    mapId: string;
	    title: string;
	    side: string;
	    notes: string;
	    markers: PlanMarker[];
	    lines: PlanLine[];
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new MapPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mapId = source["mapId"];
	        this.title = source["title"];
	        this.side = source["side"];
	        this.notes = source["notes"];
	        this.markers = this.convertValues(source["markers"], PlanMarker);
	        this.lines = this.convertValues(source["lines"], PlanLine);
	        this.updatedAt = source["updatedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace wailsiface {
	
	export class ChatMessage {
	    role: string;
	    content: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class ChatState {
	    available: boolean;
	    message: string;
	    history: ChatMessage[];
	
	    static createFrom(source: any = {}) {
	        return new ChatState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.message = source["message"];
	        this.history = this.convertValues(source["history"], ChatMessage);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LiveAnalysisResult {
	    report: analysis.Report;
	    source: string;
	    cached: boolean;
	    fetchedAt: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LiveAnalysisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.report = this.convertValues(source["report"], analysis.Report);
	        this.source = source["source"];
	        this.cached = source["cached"];
	        this.fetchedAt = source["fetchedAt"];
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LastReportResult {
	    hasReport: boolean;
	    result: LiveAnalysisResult;
	
	    static createFrom(source: any = {}) {
	        return new LastReportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hasReport = source["hasReport"];
	        this.result = this.convertValues(source["result"], LiveAnalysisResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class RiotPlayerInfo {
	    puuid: string;
	    gameName: string;
	    tagLine: string;
	    region: string;
	    shard: string;
	
	    static createFrom(source: any = {}) {
	        return new RiotPlayerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.gameName = source["gameName"];
	        this.tagLine = source["tagLine"];
	        this.region = source["region"];
	        this.shard = source["shard"];
	    }
	}
	export class RiotLoginResult {
	    success: boolean;
	    error?: string;
	    playerInfo?: RiotPlayerInfo;
	
	    static createFrom(source: any = {}) {
	        return new RiotLoginResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.playerInfo = this.convertValues(source["playerInfo"], RiotPlayerInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

