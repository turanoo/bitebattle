import SwiftUI

struct PollsView: View {
    @State private var polls: [Poll] = []
    @State private var isLoading: Bool = false
    @State private var statusMessage: String?
    @State private var showAddPoll: Bool = false
    @State private var newPollName: String = ""
    @State private var isCreatingPoll: Bool = false

    struct Poll: Identifiable, Decodable {
        let id: String
        let name: String
        let role: String
        let invite_code: String?
    }

    var body: some View {
        ZStack {
            LinearGradient(
                gradient: Gradient(colors: [Color.pink.opacity(0.7), Color.orange.opacity(0.7)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()

            VStack(spacing: 24) {
                Spacer()

                Image(systemName: "fork.knife.circle.fill")
                    .resizable()
                    .scaledToFit()
                    .frame(width: 80, height: 80)
                    .foregroundColor(.white)
                    .shadow(radius: 10)

                Text("Your Polls")
                    .font(.largeTitle)
                    .fontWeight(.bold)
                    .foregroundColor(.white)
                    .shadow(radius: 4)

                if isLoading {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                        .padding()
                } else if let status = statusMessage {
                    Text(status)
                        .foregroundColor(.white)
                        .font(.footnote)
                        .padding(.top, 8)
                        .shadow(radius: 1)
                } else if polls.isEmpty {
                    Text("No polls found.")
                        .foregroundColor(.white)
                        .font(.headline)
                        .padding()
                        .shadow(radius: 1)
                }

                ScrollView {
                    VStack(spacing: 18) {
                        ForEach(polls) { poll in
                            NavigationLink(destination: Text("Poll Details for \(poll.name)")) {
                                VStack(alignment: .leading, spacing: 8) {
                                    HStack {
                                        Text(poll.name)
                                            .font(.headline)
                                            .foregroundColor(.pink)
                                        Spacer()
                                        Text(poll.role.capitalized)
                                            .font(.caption)
                                            .foregroundColor(.white)
                                            .padding(.horizontal, 8)
                                            .padding(.vertical, 2)
                                            .background(poll.role == "owner" ? Color.pink.opacity(0.7) : Color.gray.opacity(0.5))
                                            .cornerRadius(8)
                                    }
                                    if poll.role == "owner", let code = poll.invite_code {
                                        HStack {
                                            Image(systemName: "person.2.fill")
                                                .foregroundColor(.orange)
                                            Text("Invite Code: ")
                                                .font(.subheadline)
                                                .foregroundColor(.white)
                                            Text(code)
                                                .font(.system(.subheadline, design: .monospaced))
                                                .foregroundColor(.white)
                                                .padding(.horizontal, 6)
                                                .background(Color.pink.opacity(0.5))
                                                .cornerRadius(6)
                                        }
                                    }
                                }
                                .padding()
                                .background(Color.white.opacity(0.15))
                                .cornerRadius(14)
                                .shadow(radius: 3)
                            }
                        }

                        // Add Poll Tile
                        VStack(alignment: .leading, spacing: 8) {
                            if showAddPoll {
                                TextField("Poll Name", text: $newPollName)
                                    .padding()
                                    .background(Color.white.opacity(0.9))
                                    .cornerRadius(10)
                                    .autocapitalization(.words)
                                    .disableAutocorrection(true)
                                    .shadow(radius: 1)

                                Button(action: {
                                    createPoll()
                                }) {
                                    if isCreatingPoll {
                                        ProgressView()
                                            .progressViewStyle(CircularProgressViewStyle(tint: .pink))
                                            .frame(maxWidth: .infinity)
                                            .padding()
                                            .background(Color.pink.opacity(0.8))
                                            .foregroundColor(.white)
                                            .cornerRadius(12)
                                            .shadow(radius: 2)
                                    } else {
                                        Text("Create Poll")
                                            .fontWeight(.semibold)
                                            .frame(maxWidth: .infinity)
                                            .padding()
                                            .background(newPollName.isEmpty ? Color.gray.opacity(0.5) : Color.pink.opacity(0.8))
                                            .foregroundColor(.white)
                                            .cornerRadius(12)
                                            .shadow(radius: 2)
                                    }
                                }
                                .disabled(isCreatingPoll || newPollName.isEmpty)

                                Button(action: {
                                    showAddPoll = false
                                    newPollName = ""
                                }) {
                                    Text("Cancel")
                                        .foregroundColor(.pink)
                                        .padding(.top, 4)
                                }
                            } else {
                                Button(action: {
                                    showAddPoll = true
                                }) {
                                    HStack {
                                        Image(systemName: "plus.circle.fill")
                                            .foregroundColor(.white)
                                        Text("Add New Poll")
                                            .fontWeight(.semibold)
                                            .foregroundColor(.white)
                                    }
                                    .frame(maxWidth: .infinity)
                                    .padding()
                                    .background(Color.pink.opacity(0.7))
                                    .cornerRadius(14)
                                    .shadow(radius: 2)
                                }
                            }
                        }
                        .padding()
                        .background(Color.white.opacity(0.10))
                        .cornerRadius(14)
                        .shadow(radius: 2)
                    }
                    .padding(.horizontal, 12)
                }

                Spacer()
            }
            .padding()
        }
        .navigationTitle("BiteBattle")
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                NavigationLink(destination: AccountView()) {
                    HStack(spacing: 6) {
                        Image(systemName: "person.crop.circle.fill")
                            .font(.title2)
                            .foregroundColor(.white)
                        Text("Account")
                            .foregroundColor(.white)
                            .fontWeight(.semibold)
                    }
                    .padding(.vertical, 6)
                    .padding(.horizontal, 10)
                    .background(Color.pink.opacity(0.7))
                    .cornerRadius(10)
                    .shadow(radius: 2)
                }
            }
        }
        .onAppear(perform: fetchPolls)
    }

    func fetchPolls() {
        guard let token = UserDefaults.standard.string(forKey: "authToken"),
              !token.isEmpty,
              let url = URL(string: "http://localhost:8080/api/polls") else {
            statusMessage = "Not logged in."
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        isLoading = true
        statusMessage = nil
        polls = []

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isLoading = false
                if let error = error {
                    statusMessage = "Failed: \(error.localizedDescription)"
                    return
                }
                guard let data = data else {
                    statusMessage = "No data returned."
                    return
                }
                do {
                    let decoded = try JSONDecoder().decode([Poll].self, from: data)
                    self.polls = decoded
                } catch {
                    statusMessage = "Failed to load polls."
                }
            }
        }.resume()
    }

    func createPoll() {
        guard let token = UserDefaults.standard.string(forKey: "authToken"),
              !token.isEmpty,
              let url = URL(string: "http://localhost:8080/api/polls") else {
            statusMessage = "Not logged in."
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let payload: [String: String] = [
            "name": newPollName
        ]

        guard let httpBody = try? JSONSerialization.data(withJSONObject: payload) else {
            statusMessage = "Failed to encode data."
            return
        }
        request.httpBody = httpBody

        isCreatingPoll = true
        statusMessage = nil

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isCreatingPoll = false
                if let error = error {
                    statusMessage = "Failed: \(error.localizedDescription)"
                    return
                }
                guard let data = data else {
                    statusMessage = "No data returned."
                    return
                }
                do {
                    let newPoll = try JSONDecoder().decode(Poll.self, from: data)
                    self.polls.insert(newPoll, at: 0)
                    self.newPollName = ""
                    self.showAddPoll = false
                } catch {
                    statusMessage = "Failed to create poll."
                }
            }
        }.resume()
    }
}