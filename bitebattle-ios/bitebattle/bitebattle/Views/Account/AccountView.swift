import SwiftUI

struct AccountView: View {
    @State private var name: String = ""
    @State private var email: String = ""
    @State private var password: String = ""
    @State private var statusMessage: String?
    @State private var isLoading: Bool = false

    @State private var originalName: String = ""
    @State private var originalEmail: String = ""
    @State private var fetchFailed: Bool = false
    @AppStorage("isLoggedIn") var isLoggedIn: Bool = false

    var body: some View {
        ZStack {
            LinearGradient(
                gradient: Gradient(colors: [Color.pink.opacity(0.7), Color.orange.opacity(0.7)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()

            VStack(spacing: 28) {
                Spacer()

                Image(systemName: "person.crop.circle.fill")
                    .resizable()
                    .scaledToFit()
                    .frame(width: 70, height: 70)
                    .foregroundColor(.white)
                    .shadow(radius: 8)

                Text("Account Info")
                    .font(.largeTitle)
                    .fontWeight(.bold)
                    .foregroundColor(.white)
                    .shadow(radius: 3)

                if fetchFailed {
                    Text("Error Retrieving Account Information")
                        .foregroundColor(.white)
                        .font(.headline)
                        .padding()
                        .background(Color.red.opacity(0.7))
                        .cornerRadius(10)
                        .shadow(radius: 2)
                } else {
                    VStack(spacing: 16) {
                        TextField("Name", text: $name)
                            .padding()
                            .background(Color.white.opacity(0.9))
                            .cornerRadius(10)
                            .autocapitalization(.words)
                            .disableAutocorrection(true)
                            .shadow(radius: 1)

                        TextField("Email", text: $email)
                            .padding()
                            .background(Color.white.opacity(0.9))
                            .cornerRadius(10)
                            .keyboardType(.emailAddress)
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                            .shadow(radius: 1)

                        SecureField("New Password (optional)", text: $password)
                            .padding()
                            .background(Color.white.opacity(0.9))
                            .cornerRadius(10)
                            .shadow(radius: 1)
                    }
                    .padding(.horizontal, 24)

                    Button(action: {
                        updateAccount()
                    }) {
                        if isLoading {
                            ProgressView()
                                .progressViewStyle(CircularProgressViewStyle(tint: .pink))
                                .frame(maxWidth: .infinity)
                                .padding()
                                .background(Color.pink.opacity(0.8))
                                .foregroundColor(.white)
                                .cornerRadius(12)
                                .shadow(radius: 2)
                        } else {
                            Text("Update Account")
                                .fontWeight(.semibold)
                                .frame(maxWidth: .infinity)
                                .padding()
                                .background((fieldsChanged || !password.isEmpty) ? Color.pink.opacity(0.8) : Color.gray.opacity(0.5))
                                .foregroundColor(.white)
                                .cornerRadius(12)
                                .shadow(radius: 2)
                        }
                    }
                    .padding(.horizontal, 24)
                    .disabled(isLoading || (!fieldsChanged && password.isEmpty))
                }

                if let status = statusMessage {
                    Text(status)
                        .foregroundColor(.white)
                        .font(.footnote)
                        .padding(.top, 8)
                        .shadow(radius: 1)
                }

                // Sign Out Button
                Button(action: {
                    UserDefaults.standard.removeObject(forKey: "authToken")
                    isLoggedIn = false
                }) {
                    Text("Sign Out")
                        .fontWeight(.semibold)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.red.opacity(0.8))
                        .foregroundColor(.white)
                        .cornerRadius(12)
                        .shadow(radius: 2)
                }
                .padding(.horizontal, 24)
                .padding(.top, 8)

                Spacer()
            }
            .padding()
        }
        .onAppear(perform: fetchAccount)
    }

    private var fieldsChanged: Bool {
        name != originalName || email != originalEmail
    }

    private func getAuthToken() -> String? {
        let token = UserDefaults.standard.string(forKey: "authToken")
        print("Auth Token: \(token ?? "nil")")
        return token?.isEmpty == false ? token : nil
    }

    func fetchAccount() {
        guard let token = getAuthToken(),
              let url = URL(string: "http://localhost:8080/api/account") else {
            statusMessage = "Not logged in or token missing."
            fetchFailed = true
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        isLoading = true
        fetchFailed = false
        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isLoading = false
                if let error = error {
                    statusMessage = "Fetch failed: \(error.localizedDescription)"
                    fetchFailed = true
                    return
                }
                guard let data = data else {
                    statusMessage = "No data returned."
                    fetchFailed = true
                    return
                }
                if let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
                   let fetchedName = json["name"] as? String,
                   let fetchedEmail = json["email"] as? String {
                    self.name = fetchedName
                    self.email = fetchedEmail
                    self.originalName = fetchedName
                    self.originalEmail = fetchedEmail
                    statusMessage = nil
                    fetchFailed = false
                } else {
                    fetchFailed = true
                }
            }
        }.resume()
    }

    func updateAccount() {
        guard let token = getAuthToken(),
              let url = URL(string: "http://localhost:8080/api/account") else {
            statusMessage = "Not logged in or token missing."
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        var payload: [String: String] = [
            "name": name,
            "email": email
        ]
        if !password.isEmpty {
            payload["password"] = password
        }

        guard let httpBody = try? JSONSerialization.data(withJSONObject: payload) else {
            statusMessage = "Failed to encode request body."
            return
        }
        request.httpBody = httpBody

        isLoading = true
        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isLoading = false
                if let error = error {
                    statusMessage = "Update failed: \(error.localizedDescription)"
                    return
                }
                guard let httpResponse = response as? HTTPURLResponse else {
                    statusMessage = "Invalid response"
                    return
                }

                if httpResponse.statusCode == 200 {
                    statusMessage = "Account updated successfully!"
                    password = ""
                    originalName = name
                    originalEmail = email
                } else {
                    statusMessage = "Update failed: Status \(httpResponse.statusCode)"
                    if let data = data, let responseText = String(data: data, encoding: .utf8) {
                        print("Server response: \(responseText)")
                    }
                }
            }
        }.resume()
    }
}
