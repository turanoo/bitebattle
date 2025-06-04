// Endpoints.swift
// Centralized API endpoints

struct Endpoints {
    static let base = "http://localhost:8080/api"
    static var polls: String { "\(base)/polls" }
    static func joinPoll(pollId: String) -> String { "\(base)/polls/\(pollId)/join" }
    static func pollResults(_ id: String) -> String { "\(base)/polls/\(id)/results" }
    static func pollOptions(_ id: String) -> String { "\(base)/polls/\(id)/options" }
    static var restaurantSearch: String { "\(base)/restaurants/search" }
    static var account: String { "\(base)/account" }
    static var login: String { "\(base)/login" }
    static var register: String { "\(base)/register" }
    static func updatePoll(_ id: String) -> String { "\(base)/polls/\(id)" }
    static func deletePoll(_ id: String) -> String { "\(base)/polls/\(id)" }
}
