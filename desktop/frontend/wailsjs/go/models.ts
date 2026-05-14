export namespace main {
	
	export class AppInfo {
	    name: string;
	    status: string;
	    stack: string[];
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.status = source["status"];
	        this.stack = source["stack"];
	    }
	}

}

export namespace wailsiface {
	
	export class AssistantQueryInput {
	    mapName: string;
	    agent: string;
	    side: string;
	    phase: string;
	    credits: number;
	    previousOutcome: string;
	
	    static createFrom(source: any = {}) {
	        return new AssistantQueryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mapName = source["mapName"];
	        this.agent = source["agent"];
	        this.side = source["side"];
	        this.phase = source["phase"];
	        this.credits = source["credits"];
	        this.previousOutcome = source["previousOutcome"];
	    }
	}
	export class EconomyAdviceDTO {
	    plan: string;
	    summary: string;
	    buyThreshold: number;
	    nextRoundMin: number;
	    reminder: string;
	
	    static createFrom(source: any = {}) {
	        return new EconomyAdviceDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.plan = source["plan"];
	        this.summary = source["summary"];
	        this.buyThreshold = source["buyThreshold"];
	        this.nextRoundMin = source["nextRoundMin"];
	        this.reminder = source["reminder"];
	    }
	}
	export class TacticalCardDTO {
	    id: string;
	    mapName: string;
	    agent: string;
	    side: string;
	    phase: string;
	    category: string;
	    title: string;
	    summary: string;
	    action: string;
	    priority: number;
	    safetyNotes: string;
	
	    static createFrom(source: any = {}) {
	        return new TacticalCardDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.mapName = source["mapName"];
	        this.agent = source["agent"];
	        this.side = source["side"];
	        this.phase = source["phase"];
	        this.category = source["category"];
	        this.title = source["title"];
	        this.summary = source["summary"];
	        this.action = source["action"];
	        this.priority = source["priority"];
	        this.safetyNotes = source["safetyNotes"];
	    }
	}
	export class AssistantResultDTO {
	    cards: TacticalCardDTO[];
	    economyAdvice: EconomyAdviceDTO;
	    safetyNotes: string[];
	
	    static createFrom(source: any = {}) {
	        return new AssistantResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cards = this.convertValues(source["cards"], TacticalCardDTO);
	        this.economyAdvice = this.convertValues(source["economyAdvice"], EconomyAdviceDTO);
	        this.safetyNotes = source["safetyNotes"];
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
	export class ClearCacheResult {
	    cleared: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ClearCacheResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cleared = source["cleared"];
	        this.message = source["message"];
	    }
	}
	
	export class ExportDataResult {
	    path: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ExportDataResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.message = source["message"];
	    }
	}
	export class FindingDTO {
	    type: string;
	    severity: string;
	    confidence: number;
	    title: string;
	    description: string;
	    evidence: string[];
	
	    static createFrom(source: any = {}) {
	        return new FindingDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.severity = source["severity"];
	        this.confidence = source["confidence"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.evidence = source["evidence"];
	    }
	}
	export class LookupPlayerInput {
	    name: string;
	    tag: string;
	    region: string;
	    consent: boolean;
	    apiKey: string;
	
	    static createFrom(source: any = {}) {
	        return new LookupPlayerInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.tag = source["tag"];
	        this.region = source["region"];
	        this.consent = source["consent"];
	        this.apiKey = source["apiKey"];
	    }
	}
	export class PlayerDTO {
	    puuid: string;
	    name: string;
	    tag: string;
	    region: string;
	    accountLevel: number;
	    cardSmall: string;
	    cardLarge: string;
	    lastUpdate: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.name = source["name"];
	        this.tag = source["tag"];
	        this.region = source["region"];
	        this.accountLevel = source["accountLevel"];
	        this.cardSmall = source["cardSmall"];
	        this.cardLarge = source["cardLarge"];
	        this.lastUpdate = source["lastUpdate"];
	    }
	}
	export class LookupPlayerResult {
	    player: PlayerDTO;
	    provider: string;
	    consentVersion: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LookupPlayerResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player = this.convertValues(source["player"], PlayerDTO);
	        this.provider = source["provider"];
	        this.consentVersion = source["consentVersion"];
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
	export class MatchDTO {
	    matchId: string;
	    mapName: string;
	    mode: string;
	    queue: string;
	    region: string;
	    gameStart: number;
	    roundsPlayed: number;
	    agent: string;
	    team: string;
	    kills: number;
	    deaths: number;
	    assists: number;
	    headshots: number;
	    damageMade: number;
	
	    static createFrom(source: any = {}) {
	        return new MatchDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matchId = source["matchId"];
	        this.mapName = source["mapName"];
	        this.mode = source["mode"];
	        this.queue = source["queue"];
	        this.region = source["region"];
	        this.gameStart = source["gameStart"];
	        this.roundsPlayed = source["roundsPlayed"];
	        this.agent = source["agent"];
	        this.team = source["team"];
	        this.kills = source["kills"];
	        this.deaths = source["deaths"];
	        this.assists = source["assists"];
	        this.headshots = source["headshots"];
	        this.damageMade = source["damageMade"];
	    }
	}
	
	export class RankDTO {
	    tier: number;
	    tierName: string;
	    rankingInTier: number;
	    mmrChangeToLast: number;
	    elo: number;
	    seasonId: string;
	    region: string;
	    fetchedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new RankDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tier = source["tier"];
	        this.tierName = source["tierName"];
	        this.rankingInTier = source["rankingInTier"];
	        this.mmrChangeToLast = source["mmrChangeToLast"];
	        this.elo = source["elo"];
	        this.seasonId = source["seasonId"];
	        this.region = source["region"];
	        this.fetchedAt = source["fetchedAt"];
	    }
	}
	export class RecommendationDTO {
	    title: string;
	    drill: string;
	    priority: string;
	    reason: string;
	    evidence: string[];
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new RecommendationDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.drill = source["drill"];
	        this.priority = source["priority"];
	        this.reason = source["reason"];
	        this.evidence = source["evidence"];
	        this.status = source["status"];
	    }
	}
	export class RefreshMatchesInput {
	    puuid: string;
	    region: string;
	    size: string;
	    apiKey: string;
	
	    static createFrom(source: any = {}) {
	        return new RefreshMatchesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.region = source["region"];
	        this.size = source["size"];
	        this.apiKey = source["apiKey"];
	    }
	}
	export class RefreshMatchesResult {
	    matches: MatchDTO[];
	    imported: number;
	    cached: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new RefreshMatchesResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matches = this.convertValues(source["matches"], MatchDTO);
	        this.imported = source["imported"];
	        this.cached = source["cached"];
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
	export class RefreshRankInput {
	    puuid: string;
	    region: string;
	    apiKey: string;
	
	    static createFrom(source: any = {}) {
	        return new RefreshRankInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.region = source["region"];
	        this.apiKey = source["apiKey"];
	    }
	}
	export class RefreshRankResult {
	    rank: RankDTO;
	    cached: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new RefreshRankResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rank = this.convertValues(source["rank"], RankDTO);
	        this.cached = source["cached"];
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
	export class ReportDTO {
	    id: number;
	    playerPuuid: string;
	    matchCount: number;
	    averageKda: number;
	    headshotPercent: number;
	    averageDamage: number;
	    topAgent: string;
	    topMap: string;
	    summary: string;
	    findings: FindingDTO[];
	    recommendations: RecommendationDTO[];
	
	    static createFrom(source: any = {}) {
	        return new ReportDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.playerPuuid = source["playerPuuid"];
	        this.matchCount = source["matchCount"];
	        this.averageKda = source["averageKda"];
	        this.headshotPercent = source["headshotPercent"];
	        this.averageDamage = source["averageDamage"];
	        this.topAgent = source["topAgent"];
	        this.topMap = source["topMap"];
	        this.summary = source["summary"];
	        this.findings = this.convertValues(source["findings"], FindingDTO);
	        this.recommendations = this.convertValues(source["recommendations"], RecommendationDTO);
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
	export class ResetResult {
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ResetResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message = source["message"];
	    }
	}
	export class SaveLanguageInput {
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveLanguageInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	    }
	}
	export class SaveSettingsInput {
	    apiKey: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveSettingsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiKey = source["apiKey"];
	    }
	}
	export class SettingsDTO {
	    apiKeyConfigured: boolean;
	    language: string;
	    dataPath: string;
	    cacheEntries: number;
	    expiredCacheEntries: number;
	    players: number;
	    matches: number;
	    rankSnapshots: number;
	    reports: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiKeyConfigured = source["apiKeyConfigured"];
	        this.language = source["language"];
	        this.dataPath = source["dataPath"];
	        this.cacheEntries = source["cacheEntries"];
	        this.expiredCacheEntries = source["expiredCacheEntries"];
	        this.players = source["players"];
	        this.matches = source["matches"];
	        this.rankSnapshots = source["rankSnapshots"];
	        this.reports = source["reports"];
	        this.message = source["message"];
	    }
	}

}

