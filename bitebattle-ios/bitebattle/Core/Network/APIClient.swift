// APIClient.swift
// Centralized API logic

import Foundation
import CoreLocation

// No import needed for models if in same target, just ensure files are in Compile Sources
// Endpoints.swift is in the same target, so Endpoints is available

final class APIClient {
    static let shared = APIClient()
    private init() {}

    private var token: String? {
        UserDefaults.standard.string(forKey: "authToken")
    }

    private func makeRequest(url: URL, method: String = "GET", body: Data? = nil, contentType: String? = nil) -> URLRequest? {
        guard let token = token, !token.isEmpty else { return nil }
        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        if let contentType = contentType {
            request.setValue(contentType, forHTTPHeaderField: "Content-Type")
        }
        request.httpBody = body
        return request
    }

    // MARK: - Polls
    func fetchPolls(completion: @escaping (Result<[Poll], Error>) -> Void) {
        guard let url = URL(string: Endpoints.polls), let request = makeRequest(url: url) else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { data, _, error in
            if let error = error { completion(.failure(error)); return }
            guard let data = data else { completion(.failure(APIError.noData)); return }
            do {
                let polls = try JSONDecoder().decode([Poll].self, from: data)
                completion(.success(polls))
            } catch { completion(.failure(error)) }
        }.resume()
    }

    func createPoll(name: String, completion: @escaping (Result<Poll, Error>) -> Void) {
        guard let url = URL(string: Endpoints.polls),
              let body = try? JSONSerialization.data(withJSONObject: ["name": name]),
              let request = makeRequest(url: url, method: "POST", body: body, contentType: "application/json") else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { data, _, error in
            DispatchQueue.main.async { // Ensure UI state changes are on main thread
                if let error = error { completion(.failure(error)); return }
                guard let data = data else { completion(.failure(APIError.noData)); return }
                do {
                    let poll = try JSONDecoder().decode(Poll.self, from: data)
                    completion(.success(poll))
                } catch { completion(.failure(error)) }
            }
        }.resume()
    }

    // Updated joinPoll to use pollId and inviteCode in body
    // Make sure to import or define the Poll model at the top if needed:
    // import YourModelModule

    func joinPoll(pollId: String, inviteCode: String, completion: @escaping (Result<Poll, Error>) -> Void) {
        guard let url = URL(string: Endpoints.joinPoll(pollId: pollId)),
              let body = try? JSONSerialization.data(withJSONObject: ["invite_code": inviteCode]),
              let request = makeRequest(url: url, method: "POST", body: body, contentType: "application/json") else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { data, _, error in
            if let error = error { completion(.failure(error)); return }
            guard let data = data else { completion(.failure(APIError.noData)); return }
            do {
                let poll = try JSONDecoder().decode(Poll.self, from: data)
                completion(.success(poll))
            } catch { completion(.failure(error)) }
        }.resume()
    }

    func fetchPollResults(pollId: String, completion: @escaping (Result<[PollOptionResult], Error>) -> Void) {
        guard let url = URL(string: Endpoints.pollResults(pollId)), let request = makeRequest(url: url) else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { data, _, error in
            if let error = error { completion(.failure(error)); return }
            guard let data = data else { completion(.failure(APIError.noData)); return }
            do {
                let results = try JSONDecoder().decode([PollOptionResult].self, from: data)
                completion(.success(results))
            } catch { completion(.failure(error)) }
        }.resume()
    }

    func addPollOptions(pollId: String, options: [[String: Any]], completion: @escaping (Result<Void, Error>) -> Void) {
        guard let url = URL(string: Endpoints.pollOptions(pollId)),
              let body = try? JSONSerialization.data(withJSONObject: options),
              let request = makeRequest(url: url, method: "POST", body: body, contentType: "application/json") else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { _, _, error in
            if let error = error { completion(.failure(error)); return }
            completion(.success(()))
        }.resume()
    }

    func searchRestaurants(category: String, location: String, completion: @escaping (Result<[Restaurant], Error>) -> Void) {
        guard let url = URL(string: "\(Endpoints.restaurantSearch)?q=\(category.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? "")&location=\(location.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? "")"), let request = makeRequest(url: url) else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { data, _, error in
            if let error = error { completion(.failure(error)); return }
            guard let data = data else { completion(.failure(APIError.noData)); return }
            do {
                let restaurants = try JSONDecoder().decode([Restaurant].self, from: data)
                completion(.success(restaurants))
            } catch { completion(.failure(error)) }
        }.resume()
    }

    func updatePoll(pollId: String, name: String, completion: @escaping (Result<Void, Error>) -> Void) {
        guard let url = URL(string: Endpoints.updatePoll(pollId)),
              let body = try? JSONSerialization.data(withJSONObject: ["name": name]),
              let request = makeRequest(url: url, method: "PUT", body: body, contentType: "application/json") else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { _, response, error in
            if let error = error { completion(.failure(error)); return }
            completion(.success(()))
        }.resume()
    }

    func deletePoll(pollId: String, completion: @escaping (Result<Void, Error>) -> Void) {
        guard let url = URL(string: Endpoints.deletePoll(pollId)),
              let request = makeRequest(url: url, method: "DELETE") else {
            completion(.failure(APIError.notLoggedIn)); return
        }
        URLSession.shared.dataTask(with: request) { _, response, error in
            if let error = error { completion(.failure(error)); return }
            completion(.success(()))
        }.resume()
    }

    enum APIError: Error {
        case notLoggedIn
        case noData
    }
}
